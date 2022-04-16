package dummy2

import (
	"errors"
	"log"
	"strings"

	"github.com/edsonmichaque/omni"
)

type Dummy2 struct{}

func (d Dummy2) Admit(session omni.Session, b []byte) bool {
	log.Println(session.ID, "found:", string(b))
	toUpper := strings.ToUpper(string(b))
	log.Println(session.ID, "toUpper:", toUpper)

	admit := strings.TrimSpace(toUpper) == "HELLO"

	log.Println(session.ID, "Admiting:", admit)

	return admit
}

func (d Dummy2) Decode(session omni.Session, b []byte) (*omni.Signal, error) {
	if strings.ToUpper(strings.TrimSpace(string(b))) == "POSITION UPDATE" {
		return &omni.Signal{
			PositionUpdate: &omni.PositionUpdate{},
		}, nil
	}

	return nil, errors.New("unknow signal")
}

func (d Dummy2) Encode(session omni.Session, in omni.EncodeInput) ([]byte, error) {
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

func (d Dummy2) Authorize(omni.Session, omni.Device, map[string]string) (bool, error) {
	return true, nil
}