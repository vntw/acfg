package portalloc

import (
	"testing"
)

func TestAllocPorts(t *testing.T) {
	selectableHttpPorts := []int{
		9100, 9102, 9104, 9106, 9108, 9110, 9112,
		9114, 9116, 9118, 9120, 9122, 9124, 9126,
	}
	portRange := PortRange{9000, 9100}

	pr := PortReqs{
		Tcp:  PortReq{Num: 2, Provider: portRange},
		Udp:  PortReq{Num: 3, Provider: RandomPort{}},
		Http: PortReq{Num: 4, Provider: PortSelection{selectableHttpPorts}},
	}

	sr := SpecReqs{
		Server: pr,
	}

	sm, err := AllocPorts(sr)
	if err != nil {
		t.Fatalf("error getting ports: %v", err)
	}

	if len(sm.Server.Tcp.List()) != 2 {
		t.Fatalf("unexpected tcp port length, expected 2 got %v", len(sm.Server.Tcp.List()))
	}
	if len(sm.Server.Udp.List()) != 3 {
		t.Fatalf("unexpected udp port length, expected 3 got %v", len(sm.Server.Udp.List()))
	}
	if len(sm.Server.Http.List()) != 4 {
		t.Fatalf("unexpected http port length, expected 4 got %v", len(sm.Server.Http.List()))
	}

	dups := map[int]struct{}{}
	ports := append(sm.Server.Tcp.List(), append(sm.Server.Udp.List(), sm.Server.Http.List()...)...)
	for i := 0; i < len(ports); i++ {
		if _, ok := dups[ports[i]]; ok {
			t.Fatalf("duplicate port found %d in %v", ports[i], ports)
		}
		dups[ports[i]] = struct{}{}
	}

	mappedTcpPorts := sm.Server.Tcp.List()
	for pt := 0; pt < len(mappedTcpPorts); pt++ {
		if mappedTcpPorts[pt] < portRange.Min || mappedTcpPorts[pt] > portRange.Max {
			t.Fatalf("port %d is out of %d-%d range", mappedTcpPorts[pt], portRange.Min, portRange.Max)
		}
	}

	mappedHttpPorts := sm.Server.Http.List()
HTTP:
	for pt := 0; pt < len(mappedHttpPorts); pt++ {
		for ptt := 0; ptt < len(selectableHttpPorts); ptt++ {
			if selectableHttpPorts[ptt] == mappedHttpPorts[pt] {
				continue HTTP
			}
		}

		t.Fatalf("could not validate http port %d, have %v got %v", mappedHttpPorts[pt], selectableHttpPorts, mappedHttpPorts)
	}
}
