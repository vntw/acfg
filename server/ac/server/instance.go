package server

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"acfg/ac/spec"
)

type Instance struct {
	pid     int
	cmd     *exec.Cmd
	logFile *os.File

	Spec spec.ServerSpec `json:"spec"`
	Ip   string          `json:"ip"`
	Uuid string          `json:"uuid"`
}

func NewInstance(
	pid int,
	cmd *exec.Cmd,
	logFile *os.File,
	serverSpec spec.ServerSpec,
	ip string,
	uuid string,
) *Instance {
	return &Instance{
		pid:     pid,
		cmd:     cmd,
		logFile: logFile,
		Spec:    serverSpec,
		Ip:      ip,
		Uuid:    uuid,
	}
}

func (i Instance) Pid() int {
	return i.pid
}

func (i Instance) Cmd() *exec.Cmd {
	return i.cmd
}

func (i Instance) LogFile() *os.File {
	return i.logFile
}

func (i Instance) StopProcess() error {
	if err := i.Cmd().Process.Signal(syscall.SIGTERM); err != nil {
		log.Println("error terminating instance (killing now):", err)
		return i.Cmd().Process.Kill()
	}

	return nil
}

func (i Instance) IsOnline() bool {
	resp, err := fetchServerInfo(i)

	if err != nil {
		log.Println("Online check failed:", i.Uuid)
	} else {
		log.Println("Online check succeded:", i.Uuid)
	}

	return err == nil && resp.StatusCode == 200
}

type ServerInstance struct {
	Error    string   `json:"error"`
	State    int      `json:"state"`
	Status   Status   `json:"status"`
	Instance Instance `json:"instance"`
}
