package omnitrack

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

type Provider interface {
	Handshake() error
}

func WithProvider(p Provider) DispatcherOption {
	return func(d *Dispatcher) {
		if d.providers == nil {
			d.providers = make([]Provider, 0)
		}

		d.providers = append(d.providers, p)
	}
}
