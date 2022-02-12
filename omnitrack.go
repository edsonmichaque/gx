package omnitrack

import (
	"io"
)

type Dispatcher struct {
	providers []Provider
}

type DispatcherOption func(*Dispatcher)

func NewDispatcher(opts ...DispatcherOption) *Dispatcher {
	d := Dispatcher{
		providers: make([]Provider, 0),
	}

	for _, opt := range opts {
		opt(&d)
	}

	return &d
}

type NoProviderError struct {
	err error
}

func (e NoProviderError) Error() string {
	return e.err.Error()
}

func (d Dispatcher) Dispatch(r io.Reader) (Provider, error) {
	for _, provider := range d.providers {
		if err := provider.Handshake(r); err == nil {
			return provider, nil
		}
	}

	return nil, NoProviderError{}
}

type Provider interface {
	Handshake(io.Reader) error
}

func WithProvider(p Provider) DispatcherOption {
	return func(d *Dispatcher) {
		if d.providers == nil {
			d.providers = make([]Provider, 0)
		}

		d.providers = append(d.providers, p)
	}
}
