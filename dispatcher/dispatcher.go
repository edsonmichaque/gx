package dispatcher

import (
	"bufio"
	"errors"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/edsonmichaque/omni"
	"github.com/google/uuid"
)

type Queue interface {
	Send(omni.Session, interface{}) error
	Get(omni.Session) (*omni.EncodeInput, error)
}

type Logger interface {
	Debug(...interface{})
	Debugf(...interface{})
	Warn(...interface{})
	Warnf(...interface{})
}

type Dispatcher struct {
	rw        io.ReadWriter
	providers map[string]omni.Omni
	current   omni.Omni
	queue     Queue
	logger    Logger
}

func (d Dispatcher) Dispatch() {
	sessionId := uuid.NewString()

	session := omni.Session{
		ID: sessionId,
	}

	d.logger.Debug(session.ID, "reading first bytes")
	rawBytes, err := d.read()
	if err != nil {
		d.logger.Debug(session.ID, "found an error", err)
		return
	}

	d.logger.Debug(session.ID, "finding the provider")
	for n, p := range d.providers {
		if p.Admit(session, rawBytes) {
			d.logger.Debug(session.ID, "found", n, "as a provider")
			d.current = p
			break
		}
	}

	if d.current == nil {
		d.logger.Debug(session.ID, "didnt find a provider")
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		d.logger.Debug(session.ID, "processing signals")
		defer wg.Done()
		for {
			if err := d.processSignals(session); err != nil {
				d.logger.Debug(session.ID, "err found", err)
			}
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		d.logger.Debug(session.ID, "processing commands")

		defer wg.Done()
		for {
			time.Sleep(5 * time.Second)
			if err := d.processCommands(session); err != nil && err != ErrEmptyQueue {
				d.logger.Debug(session.ID, "err found", err)
			}
		}
	}(&wg)

	wg.Wait()
}

func (d Dispatcher) read() ([]byte, error) {
	return bufio.NewReader(d.rw).ReadBytes('\n')
}

func (d Dispatcher) write(bytes []byte) (int, error) {
	return d.rw.Write(bytes)
}

func (d Dispatcher) processSignals(session omni.Session) error {
	inputBytes, err := bufio.NewReader(d.rw).ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil
		}

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
			rawBytes, err = d.current.Encode(session, omni.EncodeInput{PositionUpdateResponse: &omni.PositionUpdateResponse{}})
			if err != nil {
				return err
			}
		} else {
			rawBytes, err = d.current.Encode(session, omni.EncodeInput{PositionUpdateResponse: &omni.PositionUpdateResponse{}})
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

func (d Dispatcher) processCommands(session omni.Session) error {
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

func authorize(session omni.Session, provider omni.Omni, d omni.Device, req omni.AuthorizationRequest) ([]byte, error) {
	authz, err := provider.Authorize(session, d, req.Credentials)
	if err != nil {
		return nil, err
	}

	var rawBytes []byte
	if authz {
		rawBytes, err = provider.Encode(session, omni.EncodeInput{AuthorizationResponse: &omni.AuthorizationResponse{}})
		if err != nil {
			return nil, err
		}
	} else {
		rawBytes, err = provider.Encode(session, omni.EncodeInput{AuthorizationResponse: &omni.AuthorizationResponse{}})
		if err != nil {
			return nil, err
		}
	}

	return rawBytes, nil
}

type Option func(*Dispatcher)

type logger struct{}

func (l logger) Debug(args ...interface{}) {
	log.Println(args...)
}

func (l logger) Debugf(args ...interface{}) {
	log.Println(args...)
}

func (l logger) Warn(args ...interface{}) {
	log.Println(args...)
}

func (l logger) Warnf(args ...interface{}) {
	log.Println(args...)
}

type queue struct{}

func (q queue) Send(_ omni.Session, data interface{}) error {
	return nil
}

var ErrEmptyQueue = errors.New("empty queue")

func (q queue) Get(session omni.Session) (*omni.EncodeInput, error) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if cmd := r.Intn(3); cmd == 0 {
		log.Println(session.ID, "Authorization command")
		return &omni.EncodeInput{AuthorizationResponse: &omni.AuthorizationResponse{}}, nil
	} else if cmd == 1 {
		log.Println(session.ID, "Ignite command")
		return &omni.EncodeInput{Ignite: &omni.Ignite{}}, nil
	} else {
		log.Println(session.ID, "Nothing to send")
		return nil, ErrEmptyQueue
	}
}

func New(rw io.ReadWriter, providers map[string]omni.Omni) Dispatcher {
	return Dispatcher{
		rw:        rw,
		providers: providers,
		logger:    &logger{},
		queue:     &queue{},
	}
}
