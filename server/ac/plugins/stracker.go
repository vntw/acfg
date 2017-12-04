package plugins

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/go-ini/ini"

	"github.com/venyii/acsrvmanager/server/ac/portalloc"
)

var tcpPortRange = []int{
	50042, 50043, 54242, 54243, 60023, 60024,
	62323, 62324, 42423, 42424, 23232, 23233,
}

type Stracker struct {
	pid     int
	cmd     *exec.Cmd
	logFile *os.File

	N             string `json:"name"`
	StrackerDir   string `json:"-"`
	tmpConfigPath string
	TcpPort       int `json:"tcpPort"`
	HttpPort      int `json:"httpPort"`
}

func NewStrackerPlugin(dir string) *Stracker {
	return &Stracker{
		N:           "stracker",
		StrackerDir: dir,
	}
}

func (p Stracker) Name() string {
	return p.N
}

func (p Stracker) Pid() int {
	return p.pid
}

func (p Stracker) Cmd() *exec.Cmd {
	return p.cmd
}

func (p Stracker) LogFile() *os.File {
	return p.logFile
}

func (p *Stracker) RecvPorts(pm portalloc.PortsMap) {
	p.TcpPort = pm.Tcp.PickOne()
	p.HttpPort = pm.Http.PickOne()
}

func (p Stracker) PortReqs() portalloc.PortReqs {
	return portalloc.PortReqs{
		Tcp: portalloc.PortReq{
			Num:      1,
			Provider: portalloc.PortSelection{tcpPortRange},
		},
		Http: portalloc.PortReq{
			Num:      1,
			Provider: portalloc.RandomPort{},
		},
	}
}

func (p *Stracker) CreateTmpConfig(tmpDir string) error {
	strackerFile := filepath.Join(p.StrackerDir, "stracker.ini")
	if _, err := os.Stat(strackerFile); err != nil {
		return err
	}

	iniFile, err := ini.Load(strackerFile)
	if err != nil {
		return err
	}

	log.Printf("Setting stracker values %s, %d, %d", filepath.Join(tmpDir, "server_cfg.ini"), p.TcpPort, p.HttpPort)

	iniFile.Section("STRACKER_CONFIG").Key("ac_server_cfg_ini").SetValue(filepath.Join(tmpDir, "server_cfg.ini"))
	iniFile.Section("STRACKER_CONFIG").Key("listening_port").SetValue(strconv.Itoa(p.TcpPort))
	iniFile.Section("HTTP_CONFIG").Key("listen_port").SetValue(strconv.Itoa(p.HttpPort))

	p.tmpConfigPath = filepath.Join(tmpDir, "stracker.ini")

	if err := iniFile.SaveTo(p.tmpConfigPath); err != nil {
		return err
	}

	return nil
}

func (p *Stracker) Start(logsDir string) error {
	log.Println("Starting stracker")

	strackerLogfile, err := os.Create(filepath.Join(logsDir, "stracker.log"))

	if err != nil {
		return errors.New("error creating log file for stracker plugin ")
	}

	p.logFile = strackerLogfile

	// ./stracker --stracker_ini ../stracker.ini
	strackerExecutable := filepath.Join(p.StrackerDir, "stracker")
	strackerCmd := exec.Command(strackerExecutable, "--stracker_ini", p.tmpConfigPath)
	strackerCmd.Stdout = strackerLogfile
	strackerCmd.Stderr = strackerLogfile

	if err := strackerCmd.Start(); err != nil {
		return errors.New("error starting stracker: " + err.Error())
	}

	p.cmd = strackerCmd

	log.Println("Running stracker", strackerCmd.Process.Pid)

	return nil
}

func (p *Stracker) Stop() error {
	if err := p.Cmd().Process.Signal(syscall.SIGTERM); err != nil {
		log.Println("error terminating stracker (killing now):", err)
		return p.Cmd().Process.Kill()
	}

	return nil
}
