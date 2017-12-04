package ac

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/venyii/acsrvmanager/server/ac/config"
	"github.com/venyii/acsrvmanager/server/ac/plugins"
	"github.com/venyii/acsrvmanager/server/ac/portalloc"
	"github.com/venyii/acsrvmanager/server/ac/server"
	"github.com/venyii/acsrvmanager/server/ac/spec"
	"github.com/venyii/acsrvmanager/server/app"
)

type InstanceManager interface {
	GetInstances() []server.Instance
	GetInstance(uuid string) (server.Instance, error)

	StartInstanceWithConfig(cfg app.Config, spec *spec.ServerSpec) (server.Instance, error)
	StartInstance(cfg app.Config, spec *spec.ServerSpec) (server.Instance, error)
	StopInstance(uuid string) error
}

type MemoryInstanceManager struct {
	instances []server.Instance
	mu        sync.Mutex
}

func NewMemoryInstanceManager() *MemoryInstanceManager {
	im := &MemoryInstanceManager{}
	im.instances = make([]server.Instance, 0)
	return im
}

func (im *MemoryInstanceManager) GetInstances() []server.Instance {
	return im.instances
}

func (im *MemoryInstanceManager) GetInstance(uuid string) (server.Instance, error) {
	for _, instance := range im.instances {
		if instance.Uuid == uuid {
			// instance is removed from instances by observeProcess
			return instance, nil
		}
	}

	return server.Instance{}, errors.New("instance not found")
}

func (im *MemoryInstanceManager) StartInstanceWithConfig(cfg app.Config, srvSpec *spec.ServerSpec) (server.Instance, error) {
	sr := portalloc.NewSpecReqs(srvSpec.PortReqs)
	for _, plugin := range srvSpec.Plugins {
		sr.Plugins[plugin.Name()] = plugin.PortReqs()
	}

	sm, err := portalloc.AllocPorts(*sr)
	if err != nil {
		return server.Instance{}, errors.New("could not get ports: " + err.Error())
	}

	srvSpec.RecvPorts(sm)

	tmpsc, err := CreateTmpConfig(cfg.ServerCfgsDir, srvSpec)
	if err != nil {
		return server.Instance{}, errors.New("error creating tmp config: " + err.Error())
	}

	srvSpec.TmpConfig = tmpsc

	return im.StartInstance(cfg, srvSpec)
}

func (im *MemoryInstanceManager) StartInstance(cfg app.Config, spec *spec.ServerSpec) (server.Instance, error) {
	instanceUuid := uuid.New().String()

	logDir := filepath.Join(cfg.ServerLogsDir, instanceUuid)

	if err := os.Mkdir(logDir, 0755); err != nil {
		return server.Instance{}, errors.New("Could not create instance log dir: " + err.Error())
	}

	logfile, err := os.Create(filepath.Join(logDir, "instance.log"))

	if err != nil {
		log.Println("Error creating log file: ", err)
		return server.Instance{}, errors.New("error creating log file")
	}

	cmdString := filepath.Join(cfg.ACServerDir, cfg.ACServerBinary)
	cmd := exec.Command(cmdString, "-c", spec.TmpConfig.Path(cfg.ServerCfgsDir, config.ServerCfg), "-e", spec.TmpConfig.Path(cfg.ServerCfgsDir, config.EntryList))
	cmd.Dir = cfg.ACServerDir
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	log.Println(fmt.Sprintf("Starting instance via cmd %s (tmp: %s - config: %s)", cmdString, spec.TmpConfig.Uuid, spec.Preset.Uuid))

	if err := cmd.Start(); err != nil {
		log.Println("Error starting instance", err)
		return server.Instance{}, errors.New("error starting instance")
	}

	log.Printf("Running instance %s\n", instanceUuid)

	instance := *server.NewInstance(cmd.Process.Pid, cmd, logfile, *spec, cfg.ServerIP, instanceUuid)

	im.addInstance(instance)
	go im.observeProcess(cmd)

	if !im.isServerAvailable(instance) {
		log.Println("Server not reachable")
		im.StopInstance(instance.Uuid)

		return server.Instance{}, errors.New("server not reachable")
	}

	for _, plugin := range spec.Plugins {
		log.Println("Starting plugin:", plugin.Name())
		if err := plugin.Start(logDir); err != nil {
			im.StopInstance(instanceUuid)
			return server.Instance{}, errors.New(fmt.Sprintf("error starting plugin %s:", plugin.Name()) + err.Error())
		}

		go im.observePluginProcess(&instance, plugin)
	}

	return instance, nil
}

func (im *MemoryInstanceManager) shutdownInstance(uuid string) error {
	log.Printf("Shutting down instance %s\n", uuid)

	for _, instance := range im.instances {
		if instance.Uuid == uuid {
			for _, plugin := range instance.Spec.Plugins {
				if err := plugin.Stop(); err != nil {
					log.Printf("Error while stopping plugin %s: %v", plugin.Name(), err)
				}
			}

			return instance.StopProcess()
		}
	}

	return errors.New("instance not found")
}

func (im *MemoryInstanceManager) StopInstance(uuid string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	log.Printf("Stopping instance %s\n", uuid)

	return im.shutdownInstance(uuid)
}

func (im *MemoryInstanceManager) observeProcess(cmd *exec.Cmd) {
	if err := cmd.Wait(); err != nil {
		exitErr, ok := err.(*exec.ExitError)

		if !ok {
			log.Println("Error when instance stopped", err)
		} else {
			log.Println("Instance stopped: " + exitErr.ProcessState.String())
		}
	}

	log.Println("Done observing instance", cmd.ProcessState.String())

	for i, instance := range im.instances {
		if instance.Pid() == cmd.Process.Pid {
			if err := instance.LogFile().Close(); err != nil {
				log.Println("Error closing log file")
			}

			im.deleteInstance(i)
			log.Println("Instance removed " + instance.Uuid)
			return
		}
	}
}

func (im *MemoryInstanceManager) observePluginProcess(i *server.Instance, p plugins.Plugin) {
	if err := p.Cmd().Wait(); err != nil {
		exitErr, ok := err.(*exec.ExitError)

		if !ok {
			log.Printf("error when plugin %s stopped: %v\n", p.Name(), err)
		} else {
			log.Printf("plugin %s stopped: %s\n", p.Name(), exitErr.ProcessState.String())
		}

		if err := p.LogFile().Close(); err != nil {
			log.Printf("error closing plugin %s log file\n", p.Name())
		}

		im.shutdownInstance(i.Uuid)
	}
}

// Give the server some time to start properly
func (im *MemoryInstanceManager) isServerAvailable(instance server.Instance) bool {
	for i := 0; i < 10; i++ {
		if instance.IsOnline() {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}

	return false
}

func (im *MemoryInstanceManager) addInstance(instance server.Instance) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.instances = append(im.instances, instance)
}

func (im *MemoryInstanceManager) deleteInstance(idx int) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.instances = append(im.instances[:idx], im.instances[idx+1:]...)
}
