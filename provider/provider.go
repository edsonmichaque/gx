package provider

import "io"

type Provider interface {
	Handshake(io.Reader) error
	Handle() error
}
