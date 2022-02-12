package dispatcher

import (
	"io"

	"github.com/edsonmichaque/omni/provider"
)

type Dispatcher struct {
	providers []provider.Provider
}

type Option func(*Dispatcher)

func New(opts ...Option) *Dispatcher {
	d := Dispatcher{
		providers: make([]provider.Provider, 0),
	}

	for _, opt := range opts {
		opt(&d)
	}

	return &d
}

func (d *Dispatcher) Dispatch(r io.Reader) (provider.Provider, error) {
	return d.dispatch(r)
}

func (d *Dispatcher) dispatch(r io.Reader) (provider.Provider, error) {
	for _, provider := range d.providers {
		if err := provider.Handshake(r); err == nil {
			return provider, nil
		}
	}

	return nil, MissingProviderError{}
}

func WithProvider(p provider.Provider) Option {
	return func(d *Dispatcher) {
		if d.providers == nil {
			d.providers = make([]provider.Provider, 0)
		}

		d.providers = append(d.providers, p)
	}
}

func (d *Dispatcher) Handle(r io.ReadWriteCloser) error {
	return d.handle(r)
}

func (d *Dispatcher) handle(r io.ReadWriteCloser) error {
	provider, err := d.Dispatch(r)
	if err != nil {
		return err
	}

	if err := provider.Handle(); err != nil {
		return err
	}

	return nil
}
