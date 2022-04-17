package dispatcher

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"syscall"
	"time"

	"github.com/edsonmichaque/omni/internal/logger"
	"github.com/edsonmichaque/omni/internal/logger/stdlib"
	"github.com/edsonmichaque/omni/internal/queue"
	"github.com/edsonmichaque/omni/internal/queue/memory"
	"github.com/edsonmichaque/omni/libomni"
)

type connectionState uint

const (
	connectionOpen   = iota
	connectionClosed = iota
)

type Dispatcher struct {
	state     connectionState
	rw        io.ReadWriter
	providers map[string]libomni.Omni
	current   libomni.Omni
	queue     queue.Queue
	logger    logger.Logger
}

func (d Dispatcher) Dispatch(session libomni.Session) error {

	d.logger.Info(session.ID, "reading first bytes")
	rawBytes, err := d.read()
	if err != nil {
		d.logger.Error(session.ID, "found an error", err)
		return err
	}

	d.logger.Info(session.ID, "finding the provider")
	for n, p := range d.providers {
		if p.Admit(session, rawBytes) {
			d.logger.Info(session.ID, "found", n, "as a provider")
			d.current = p
			break
		}
	}

	if d.current == nil {
		d.logger.Error(session.ID, "didnt find a provider")
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		d.logger.Info(session.ID, "processing signals")
		defer wg.Done()
		for {
			if d.state == connectionClosed {
				d.logger.Error(session.ID, "connection was closed by the client")
				d.logger.Error(session.ID, "commands processor is aborting")

				break
			}

			if err := d.processSignals(session); err != nil {
				d.logger.Info(session.ID, "err found", err)
				if err == io.EOF {
					d.state = connectionClosed

					d.logger.Error(session.ID, "connection was closed by the client")
					d.logger.Error(session.ID, "signals processor is aborting")
					break
				}
			}
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		d.logger.Info(session.ID, "processing commands")

		defer wg.Done()
		for {
			if d.state == connectionClosed {
				d.logger.Error(session.ID, "connection was closed by the client")
				d.logger.Error(session.ID, "commands processor is aborting")

				break
			}

			time.Sleep(5 * time.Second)
			if err := d.processCommands(session); err != nil && err != queue.ErrEmpty {
				d.logger.Error(session.ID, "err found", err)

				if errors.Is(err, syscall.EPIPE) {
					d.state = connectionClosed

					d.logger.Error(session.ID, "connection was closed by the client")
					d.logger.Error(session.ID, "commands processor is aborting")
					break
				}
			}
		}
	}(&wg)

	wg.Wait()

	return nil
}

func (d Dispatcher) read() ([]byte, error) {
	return bufio.NewReader(d.rw).ReadBytes('\n')
}

func (d Dispatcher) write(bytes []byte) (int, error) {
	return d.rw.Write(bytes)
}

func (d Dispatcher) processSignals(session libomni.Session) error {
	inputBytes, err := bufio.NewReader(d.rw).ReadBytes('\n')
	if err != nil {
		return err
	}

	signal, err := d.current.Decode(session, inputBytes)
	if err != nil {
		return err
	}

	switch {
	case signal.AuthorizationRequest != nil:

		rawBytes, err := authorize(session, d.current, *signal.Device, *signal.AuthorizationRequest)
		if err != nil {
			return err
		}

		if _, err := d.write(rawBytes); err != nil {
			return err
		}

	case signal.PositionUpdate != nil:

		var rawBytes []byte
		if err := d.queue.Send(session, signal.PositionUpdate); err != nil {
			rawBytes, err = d.current.Encode(session, libomni.EncodeInput{PositionUpdateResponse: &libomni.PositionUpdateResponse{}})
			if err != nil {
				return err
			}
		} else {
			rawBytes, err = d.current.Encode(session, libomni.EncodeInput{PositionUpdateResponse: &libomni.PositionUpdateResponse{}})
			if err != nil {
				return err
			}
		}

		if _, err := d.write(rawBytes); err != nil {
			return err
		}
	}

	if signal.Close != nil && *signal.Close {
		return errors.New("eof")
	}

	return errors.New("")
}

func (d Dispatcher) processCommands(session libomni.Session) error {
	cmd, err := d.queue.Get(session)
	if err != nil {
		return err
	}

	rawBytes, err := d.current.Encode(session, *cmd)
	if err != nil {
		return err
	}

	if _, err := d.rw.Write(rawBytes); err != nil {
		return err
	}

	return nil
}

func authorize(session libomni.Session, provider libomni.Omni, d libomni.Device, req libomni.AuthorizationRequest) ([]byte, error) {
	authz, err := provider.Authorize(session, d, req.Credentials)
	if err != nil {
		return nil, err
	}

	var rawBytes []byte
	if authz {
		rawBytes, err = provider.Encode(session, libomni.EncodeInput{AuthorizationResponse: &libomni.AuthorizationResponse{}})
		if err != nil {
			return nil, err
		}
	} else {
		rawBytes, err = provider.Encode(session, libomni.EncodeInput{AuthorizationResponse: &libomni.AuthorizationResponse{}})
		if err != nil {
			return nil, err
		}
	}

	return rawBytes, nil
}

func New(rw io.ReadWriter, providers map[string]libomni.Omni) Dispatcher {
	return Dispatcher{
		rw:        rw,
		providers: providers,
		logger:    &stdlib.Logger{},
		queue:     &memory.Queue{},
	}
}
