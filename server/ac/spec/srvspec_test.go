package spec

import (
	"testing"

	"acfg/ac/plugins"
	"acfg/ac/portalloc"
)

type dummyPlugin struct {
	plugins.Plugin

	PM portalloc.PortsMap
}

func (dp *dummyPlugin) Name() string {
	return "dummy"
}

func (dp *dummyPlugin) RecvPorts(pm portalloc.PortsMap) {
	dp.PM = pm
}

func (dp dummyPlugin) PortReqs() portalloc.PortReqs {
	return portalloc.PortReqs{
		Tcp:  portalloc.PortReq{Num: 1},
		Udp:  portalloc.PortReq{Num: 2},
		Http: portalloc.PortReq{Num: 3},
	}
}

func TestRecvPorts(t *testing.T) {
	dp := &dummyPlugin{}

	sm := portalloc.SpecMaps{
		Server: portalloc.PortsMap{
			Tcp:  *portalloc.NewPortList([]int{1}),
			Udp:  *portalloc.NewPortList([]int{2}),
			Http: *portalloc.NewPortList([]int{3}),
		},
		Plugins: map[string]portalloc.PortsMap{
			dp.Name(): {
				Tcp:  *portalloc.NewPortList([]int{4}),
				Udp:  *portalloc.NewPortList([]int{5, 6}),
				Http: *portalloc.NewPortList([]int{7, 8, 9}),
			},
		},
	}

	srvSpec := NewServerSpec()
	srvSpec.AddPlugin(dp)
	srvSpec.RecvPorts(sm)

	if srvSpec.Ports.TcpPort != 1 {
		t.Fatalf("Unexpected tcp port, expected 1 got %v", srvSpec.Ports.TcpPort)
	}
	if srvSpec.Ports.UdpPort != 2 {
		t.Fatalf("Unexpected udp port, expected 2 got %v", srvSpec.Ports.UdpPort)
	}
	if srvSpec.Ports.HttpPort != 3 {
		t.Fatalf("Unexpected http port, expected 3 got %v", srvSpec.Ports.HttpPort)
	}

	tcpList := dp.PM.Tcp.List()
	udpList := dp.PM.Udp.List()
	httpList := dp.PM.Http.List()

	if len(tcpList) != 1 || tcpList[0] != 4 {
		t.Fatalf("Unexpected plugin tcp port slice, expected [4] got %v", tcpList)
	}
	if len(udpList) != 2 || udpList[0] != 5 || udpList[1] != 6 {
		t.Fatalf("Unexpected plugin udp port slice, expected [5,6] got %v", udpList)
	}
	if len(httpList) != 3 || httpList[0] != 7 || httpList[1] != 8 || httpList[2] != 9 {
		t.Fatalf("Unexpected plugin http port slice, expected [7,8,9] got %v", httpList)
	}
}
