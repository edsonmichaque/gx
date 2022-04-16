package registry

import (
	"github.com/edsonmichaque/omni"
	"github.com/edsonmichaque/omni/providers/tcp/dummy"
	"github.com/edsonmichaque/omni/providers/tcp/dummy2"
)

func TCPProviders() map[string]omni.Omni {
	return map[string]omni.Omni{
		"dummy":  dummy.Dummy{},
		"dummy2": dummy2.Dummy2{},
	}
}

func UDPProviders() map[string]omni.Omni {
	return map[string]omni.Omni{
		"dummy":  dummy.Dummy{},
		"dummy2": dummy2.Dummy2{},
	}
}
