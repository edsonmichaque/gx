package dummy

import (
	"log"
	"strings"

	"github.com/edsonmichaque/omni/libomni"
)

type admiter struct{}

func (d admiter) Admit(session libomni.Session, b []byte) bool {
	log.Println(session.ID, "found:", string(b))
	toUpper := strings.ToUpper(string(b))
	log.Println(session.ID, "toUpper:", toUpper)

	admit := strings.TrimSpace(toUpper) == "ADMIT"

	log.Println(session.ID, "Admiting:", admit)

	return admit
}
