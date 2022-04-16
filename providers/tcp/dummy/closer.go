package dummy

import (
	"errors"

	"github.com/edsonmichaque/omni"
)

type closer struct{}

func (d closer) Close(omni.Session) (*bool, error) {
	return nil, errors.New("not implemented")
}
