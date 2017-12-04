package ac

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"

	"github.com/venyii/acfg/server/ac/config"
	"github.com/venyii/acfg/server/ac/spec"
)

type ConfigManager interface {
	DeleteServerConfig(uuid string) error
	GetServerConfig(uuid string) (config.ServerConfigs, error)
	GetServerConfigs() []config.ServerConfigs
	CreateServerConfigs(cfgs config.ServerConfigFiles, add bool) (config.ServerConfigs, error)
	GetServerConfigByChecksums(cfgs config.ServerConfigFiles) (config.ServerConfigs, error)
}

type MemoryConfigManager struct {
	configs []config.ServerConfigs
	mu      sync.Mutex
}

func NewMemoryConfigManager() *MemoryConfigManager {
	cm := &MemoryConfigManager{}
	cm.configs = make([]config.ServerConfigs, 0)
	return cm
}

func (cm *MemoryConfigManager) addConfig(cs config.ServerConfigs) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.configs = append(cm.configs, cs)
}

func (cm *MemoryConfigManager) deleteConfig(idx int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.configs = append(cm.configs[:idx], cm.configs[idx+1:]...)
}

func (cm *MemoryConfigManager) DeleteServerConfig(uuid string) error {
	for i, cfg := range cm.configs {
		if cfg.Uuid == uuid {
			cm.deleteConfig(i)

			return nil
		}
	}

	return errors.New("no config found")
}

func (cm *MemoryConfigManager) GetServerConfig(uuid string) (config.ServerConfigs, error) {
	for _, cfg := range cm.configs {
		if cfg.Uuid == uuid {
			return cfg, nil
		}
	}

	return config.ServerConfigs{}, errors.New("no config found")
}

func (cm *MemoryConfigManager) GetServerConfigs() []config.ServerConfigs {
	return cm.configs
}

func (cm *MemoryConfigManager) GetServerConfigByChecksums(cfgs config.ServerConfigFiles) (config.ServerConfigs, error) {
	for _, sc := range cm.GetServerConfigs() {
		if sc.Files[config.ServerCfg].Checksum == cfgs[config.ServerCfg].Checksum &&
			sc.Files[config.EntryList].Checksum == cfgs[config.EntryList].Checksum {
			return sc, nil
		}
	}

	return config.ServerConfigs{}, errors.New("no config found")
}

func (cm *MemoryConfigManager) CreateServerConfigs(cfgs config.ServerConfigFiles, add bool) (config.ServerConfigs, error) {
	cs, err := config.NewServerConfigs(cfgs)
	if err != nil {
		return config.ServerConfigs{}, err
	}

	if add {
		if _, err = cm.GetServerConfigByChecksums(cfgs); err == nil {
			return config.ServerConfigs{}, errors.New("config already uploaded")
		}

		cm.addConfig(cs)
	}

	return cs, nil
}

func CreateTmpConfig(tmpDir string, srvSpec *spec.ServerSpec) (config.ServerConfigs, error) {
	tmpConfigUuid := uuid.New().String()
	tmpPath := fmt.Sprintf("%s/%s", tmpDir, tmpConfigUuid)

	if err := os.Mkdir(tmpPath, 0755); err != nil {
		return config.ServerConfigs{}, errors.New("could not create tmp cfg dir: " + err.Error())
	}

	tmpServerCfgPath := fmt.Sprintf("%s/%s", tmpPath, srvSpec.Preset.Files[config.ServerCfg].Filename())
	tmpEntryListPath := fmt.Sprintf("%s/%s", tmpPath, srvSpec.Preset.Files[config.EntryList].Filename())

	newTmpCfg := srvSpec.Preset
	newTmpCfg.Files[config.ServerCfg].Ini.Section("SERVER").Key("TCP_PORT").SetValue(strconv.Itoa(srvSpec.Ports.TcpPort))
	newTmpCfg.Files[config.ServerCfg].Ini.Section("SERVER").Key("HTTP_PORT").SetValue(strconv.Itoa(srvSpec.Ports.HttpPort))
	newTmpCfg.Files[config.ServerCfg].Ini.Section("SERVER").Key("UDP_PORT").SetValue(strconv.Itoa(srvSpec.Ports.UdpPort))

	if err := newTmpCfg.Files[config.ServerCfg].Ini.SaveTo(tmpServerCfgPath); err != nil {
		return config.ServerConfigs{}, errors.New("could not write tmp server config: " + err.Error())
	}
	if err := newTmpCfg.Files[config.EntryList].Ini.SaveTo(tmpEntryListPath); err != nil {
		return config.ServerConfigs{}, errors.New("could not write tmp entry list: " + err.Error())
	}

	for _, plugin := range srvSpec.Plugins {
		if err := plugin.CreateTmpConfig(tmpPath); err != nil {
			return config.ServerConfigs{}, errors.New("could not write tmp entry list: " + err.Error())
		}
	}

	newTmpCfg.Uuid = tmpConfigUuid

	return newTmpCfg, nil
}
