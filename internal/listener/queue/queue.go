package queue

import (
	"errors"

	"github.com/edsonmichaque/libomni"
)

type Queue interface {
	Enqueue(libomni.Session, interface{}) error
	Dequeue(libomni.Session) (*libomni.EncodeInput, error)
}

var (
	ErrEmpty = errors.New("empty")
)
