package dummy

import (
	"errors"

	"github.com/edsonmichaque/omni/libomni"
)

type closer struct{}

func (d closer) Close(libomni.Session) (*bool, error) {
	return nil, errors.New("not implemented")
}
