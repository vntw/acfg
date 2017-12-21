package spec

import (
	"github.com/venyii/acfg/server/ac/config"
	"github.com/venyii/acfg/server/ac/plugins"
	"github.com/venyii/acfg/server/ac/portalloc"
)

type ServerSpec struct {
	Preset    config.ServerConfigs  `json:"preset"`
	TmpConfig config.ServerConfigs  `json:"tmpConfig"`
	Plugins   plugins.Plugins       `json:"plugins"`
	Ports     portalloc.ServerPorts `json:"ports"`
	PortReqs  portalloc.PortReqs    `json:"-"`
}

func NewServerSpec() *ServerSpec {
	return &ServerSpec{
		// default AC server
		PortReqs: portalloc.PortReqs{
			Tcp:  portalloc.PortReq{Num: 1, Provider: portalloc.RandomPort{}},
			Udp:  portalloc.PortReq{Num: 1, Provider: portalloc.RandomPort{}},
			Http: portalloc.PortReq{Num: 1, Provider: portalloc.PortRange{Min: 8000, Max: 10000}},
		},
	}
}

func (s *ServerSpec) AddPlugin(p plugins.Plugin) {
	s.Plugins = append(s.Plugins, p)
}

func (s *ServerSpec) RecvPorts(sm portalloc.SpecMaps) {
	s.Ports.TcpPort = sm.Server.Tcp.PickOne()
	s.Ports.HttpPort = sm.Server.Http.PickOne()
	s.Ports.UdpPort = sm.Server.Udp.PickOne()

	for _, plugin := range s.Plugins {
		pm := sm.Plugins[plugin.Name()]
		plugin.RecvPorts(portalloc.PortsMap{
			Tcp:  pm.Tcp.Pick(plugin.PortReqs().Tcp.Num),
			Udp:  pm.Udp.Pick(plugin.PortReqs().Udp.Num),
			Http: pm.Http.Pick(plugin.PortReqs().Http.Num),
		})
	}
}
