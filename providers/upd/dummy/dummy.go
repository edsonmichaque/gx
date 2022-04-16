package dummy

import (
	"errors"
	"log"
	"strings"

	"github.com/edsonmichaque/omni/libomni"
)

type Dummy struct{}

func (d Dummy) Admit(sessionId string, b []byte) bool {
	log.Println(sessionId, "found:", string(b))
	toUpper := strings.ToUpper(string(b))
	log.Println(sessionId, "toUpper:", toUpper)

	admit := strings.TrimSpace(toUpper) == "ADMIT"

	log.Println(sessionId, "Admiting:", admit)

	return admit
}

func (d Dummy) Decode(sessionId string, b []byte) (*libomni.Signal, error) {
	if strings.ToUpper(strings.TrimSpace(string(b))) == "POSITION UPDATE" {
		return &libomni.Signal{
			PositionUpdate: &libomni.PositionUpdate{},
		}, nil
	}

	return nil, errors.New("unknow signal")
}

func (d Dummy) Encode(sessionId string, in libomni.EncodeInput) ([]byte, error) {
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

func (d Dummy) Authorize(string, libomni.Device, map[string]string) (bool, error) {
	return true, nil
}

func (d Dummy) Close(string) (*bool, error) {
	return nil, errors.New("not implemented")
}
