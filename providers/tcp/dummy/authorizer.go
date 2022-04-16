package dummy

import "github.com/edsonmichaque/omni"

type authorizer struct{}

func (d authorizer) Authorize(omni.Session, omni.Device, map[string]string) (bool, error) {
	return true, nil
}
