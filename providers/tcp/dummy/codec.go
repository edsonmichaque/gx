package dummy

import (
	"errors"
	"strings"

	"github.com/edsonmichaque/omni"
)

type codec struct{}

func (d codec) Decode(session omni.Session, b []byte) (*omni.Signal, error) {
	if strings.ToUpper(strings.TrimSpace(string(b))) == "POSITION UPDATE" {
		return &omni.Signal{PositionUpdate: &omni.PositionUpdate{}}, nil
	}

	return nil, errors.New("unknow signal")
}

func (d codec) Encode(session omni.Session, in omni.EncodeInput) ([]byte, error) {
	if in.AuthorizationResponse != nil {
		return []byte("authorized\n"), nil
	}

	if in.PositionUpdateResponse != nil {
		return []byte("position updated\n"), nil
	}

	if in.Ignite != nil {
		return []byte("ignite\n"), nil
	}

	return nil, errors.New("unknown command")
}
