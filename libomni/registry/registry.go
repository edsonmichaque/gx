package registry

import (
	"github.com/edsonmichaque/omni/libomni"
	"github.com/edsonmichaque/omni/libomni/providers/tcp/dummy"
	"github.com/edsonmichaque/omni/libomni/providers/tcp/dummy2"
)

func TCPProviders() map[string]libomni.Omni {
	return map[string]libomni.Omni{
		"dummy":  dummy.Provider{},
		"dummy2": dummy2.Dummy2{},
	}
}

func UDPProviders() map[string]libomni.Omni {
	return map[string]libomni.Omni{
		"dummy":  dummy.Provider{},
		"dummy2": dummy2.Dummy2{},
	}
}
