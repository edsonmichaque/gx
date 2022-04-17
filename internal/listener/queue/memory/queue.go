package memory

import (
	"log"
	"math/rand"
	"time"

	"github.com/edsonmichaque/libomni"
	"github.com/edsonmichaque/omni/internal/listener/queue"
)

func (q Queue) Get(session libomni.Session) (*libomni.EncodeInput, error) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if cmd := r.Intn(3); cmd == 0 {
		log.Println(session.ID, "Authorization command")
		return &libomni.EncodeInput{AuthorizationResponse: &libomni.AuthorizationResponse{}}, nil
	} else if cmd == 1 {
		log.Println(session.ID, "Ignite command")
		return &libomni.EncodeInput{Ignite: &libomni.Ignite{}}, nil
	} else {
		log.Println(session.ID, "Nothing to send")
		return nil, queue.ErrEmpty
	}
}

type Queue struct{}

func (q Queue) Send(_ libomni.Session, data interface{}) error {
	return nil
}
