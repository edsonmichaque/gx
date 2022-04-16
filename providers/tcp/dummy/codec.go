package dummy

import (
	"errors"
	"strings"

	"github.com/edsonmichaque/omni/libomni"
)

type codec struct{}

func (d codec) Decode(session libomni.Session, b []byte) (*libomni.Signal, error) {
	if strings.ToUpper(strings.TrimSpace(string(b))) == "POSITION UPDATE" {
		return &libomni.Signal{PositionUpdate: &libomni.PositionUpdate{}}, nil
	}

	return nil, errors.New("unknow signal")
}

func (d codec) Encode(session libomni.Session, in libomni.EncodeInput) ([]byte, error) {
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
