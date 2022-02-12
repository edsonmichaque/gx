package provider

import "io"

type Provider interface {
	Handshake(io.Reader) error
	ReplyHandshake(io.Writer) error
	Handle() error

	PositioUpdateRequest(Position) error
	PositioUpdateReply(Position) error
}

type Position struct {
	Lat  float64
	Long float64
}

type Message int

const (
	PositionUpdate Message = iota
)
