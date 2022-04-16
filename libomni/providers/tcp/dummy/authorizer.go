package dummy

import "github.com/edsonmichaque/omni/libomni"

type authorizer struct{}

func (d authorizer) Authorize(libomni.Session, libomni.Device, map[string]string) (bool, error) {
	return true, nil
}
