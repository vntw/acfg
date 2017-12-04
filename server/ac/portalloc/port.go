package portalloc

import (
	"errors"
	"io"
	"net"
)

type portProvider interface {
	get() []int
}

type RandomPort struct{}

func (rp RandomPort) get() []int {
	return []int{0}
}

type PortSelection struct {
	P []int
}

func (ps PortSelection) get() []int {
	return ps.P
}

type PortRange struct {
	Min int
	Max int
}

func (pr PortRange) get() []int {
	r := []int{}
	for i := pr.Min; i <= pr.Max; i++ {
		r = append(r, i)
	}
	return r
}

type PortReq struct {
	Num      int
	Provider portProvider
}

type PortReqs struct {
	Tcp  PortReq
	Http PortReq
	Udp  PortReq
}

func (pr *PortReqs) Sum() int {
	return pr.Tcp.Num + pr.Udp.Num + pr.Http.Num
}

type PortList struct {
	l []int
}

func NewPortList(ports []int) *PortList {
	return &PortList{l: ports}
}

func (pl *PortList) add(port int) {
	pl.l = append(pl.l, port)
}

func (pl *PortList) List() []int {
	return pl.l
}

func (pl *PortList) PickOne() int {
	return pl.Pick(1).l[0]
}

func (pl *PortList) Pick(num int) PortList {
	ext := PortList{l: pl.l[:num]}
	pl.l = pl.l[num:]

	return ext
}

type PortsMap struct {
	Tcp  PortList
	Http PortList
	Udp  PortList
}

type ServerPorts struct {
	TcpPort  int `json:"tcpPort"`
	UdpPort  int `json:"udpPort"`
	HttpPort int `json:"httpPort"`
}

type SpecReqs struct {
	Server  PortReqs
	Plugins map[string]PortReqs
}

func NewSpecReqs(serverReqs PortReqs) *SpecReqs {
	return &SpecReqs{
		Server:  serverReqs,
		Plugins: map[string]PortReqs{},
	}
}

func (sr SpecReqs) Sum() int {
	c := sr.Server.Sum()
	for _, pr := range sr.Plugins {
		c += pr.Sum()
	}
	return c
}

type SpecMaps struct {
	Server  PortsMap
	Plugins map[string]PortsMap
}

func AllocPorts(sr SpecReqs) (SpecMaps, error) {
	blocked := make(chan io.Closer, sr.Sum())

	// close blocked listeners when done
	defer func(blocked chan io.Closer) {
		close(blocked)
		for c := range blocked {
			c.Close()
		}
	}(blocked)

	spm, err := satisfyReqs(blocked, sr.Server)
	if err != nil {
		return SpecMaps{}, err
	}

	sm := SpecMaps{
		Server:  spm,
		Plugins: map[string]PortsMap{},
	}
	sm.Server = spm

	for k, pr := range sr.Plugins {
		ppm, err := satisfyReqs(blocked, pr)
		if err != nil {
			return SpecMaps{}, err
		}
		sm.Plugins[k] = ppm
	}

	return sm, nil
}

func satisfyReqs(blocked chan<- io.Closer, pr PortReqs) (PortsMap, error) {
	pm := PortsMap{}

	for i := 0; i < pr.Tcp.Num; i++ {
		port, err := reserveTcpPort(blocked, pr.Tcp.Provider)
		if err != nil {
			return PortsMap{}, errors.New("could not get tcp port")
		}
		pm.Tcp.add(port)
	}

	for i := 0; i < pr.Udp.Num; i++ {
		port, err := reserveUdpPort(blocked, pr.Udp.Provider)
		if err != nil {
			return PortsMap{}, errors.New("could not get udp port")
		}
		pm.Udp.add(port)
	}

	for i := 0; i < pr.Http.Num; i++ {
		port, err := reserveHttpPort(blocked, pr.Http.Provider)
		if err != nil {
			return PortsMap{}, errors.New("could not get http port")
		}
		pm.Http.add(port)
	}

	return pm, nil
}

func reserveTcpPort(blocked chan<- io.Closer, pr portProvider) (int, error) {
	ports := pr.get()

	if ports[0] == 0 {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return 0, errors.New("could not get tcp 0 port: " + err.Error())
		}

		ports[0] = addr.Port
	} else {
		ports = pr.get()
	}

	var err error
	for i := 0; i < len(ports); i++ {
		l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: ports[i]})
		if err != nil {
			continue
		}

		blocked <- l

		return l.Addr().(*net.TCPAddr).Port, nil
	}

	return 0, errors.New("could not get tcp port, last error: " + err.Error())
}

func reserveUdpPort(blocked chan<- io.Closer, pr portProvider) (int, error) {
	ports := pr.get()

	if ports[0] == 0 {
		addr, err := net.ResolveUDPAddr("udp", "localhost:0")
		if err != nil {
			return 0, errors.New("could not get udp 0 port: " + err.Error())
		}

		ports[0] = addr.Port
	}

	var err error
	for i := 0; i < len(ports); i++ {
		l, err := net.ListenUDP("udp", &net.UDPAddr{Port: ports[i]})
		if err != nil {
			return 0, err
		}

		blocked <- l

		return l.LocalAddr().(*net.UDPAddr).Port, nil
	}

	return 0, errors.New("could not get udp port, last error: " + err.Error())
}

func reserveHttpPort(blocked chan<- io.Closer, pr portProvider) (int, error) {
	p, err := reserveTcpPort(blocked, pr)
	if err != nil {
		return p, errors.New("could not get http port: " + err.Error())
	}

	return p, nil
}
