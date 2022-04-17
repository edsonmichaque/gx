package queue

import (
	"errors"

	"github.com/edsonmichaque/omni/libomni"
)

type Queue interface {
	Send(libomni.Session, interface{}) error
	Get(libomni.Session) (*libomni.EncodeInput, error)
}

var (
	ErrEmpty = errors.New("empty")
)
