package ac

import (
	"testing"

	"acfg/ac/config"
)

func TestCreateServerConfigs(t *testing.T) {
	srvCfg, err := config.NewServerConfig(config.ServerCfg, "[SERVER]\nNAME=Test")
	if err != nil {
		t.Fatal(err)
	}
	entryList, err := config.NewServerConfig(config.EntryList, "[CAR_0]\nMODEL=rss_formula_hybrid_2017")
	if err != nil {
		t.Fatal(err)
	}

	cfgs := config.ServerConfigFiles{
		config.ServerCfg: *srvCfg,
		config.EntryList: *entryList,
	}
	sc, err := config.NewServerConfigs(cfgs)
	if err != nil {
		t.Error("Unexpected error")
	}
	if sc.Name != "Test" {
		t.Error("Invalid name")
	}
	if sc.Uuid == "" {
		t.Error("UUID is empty")
	}
	if _, ok := sc.Files[config.ServerCfg]; !ok {
		t.Error("No server_cfg")
	}
	if _, ok := sc.Files[config.EntryList]; !ok {
		t.Error("No entry_list")
	}
}

func TestCreateServerConfigsValidatesFiles(t *testing.T) {
	var cfgs config.ServerConfigFiles
	var err error

	cfgs = config.ServerConfigFiles{}
	_, err = config.NewServerConfigs(cfgs)
	if err == nil || err.Error() != "server_cfg missing" {
		t.Error("Expected validation error")
	}

	cfgs = config.ServerConfigFiles{
		config.ServerCfg: config.ServerConfig{},
	}
	_, err = config.NewServerConfigs(cfgs)
	if err == nil || err.Error() != "entry_list missing" {
		t.Error("Expected validation error")
	}

	cfgs = config.ServerConfigFiles{
		config.ServerCfg: config.ServerConfig{},
		config.EntryList: config.ServerConfig{},
	}
	_, err = config.NewServerConfigs(cfgs)
	if err == nil || err.Error() != "server_cfg checksum missing" {
		t.Error("Expected validation error")
	}

	srvCfg, err := config.NewServerConfig(config.ServerCfg, "[SERVER]\nNAME=Test")
	if err != nil {
		t.Fatal(err)
	}
	cfgs = config.ServerConfigFiles{
		config.ServerCfg: *srvCfg,
		config.EntryList: config.ServerConfig{},
	}
	_, err = config.NewServerConfigs(cfgs)
	if err == nil || err.Error() != "entry_list checksum missing" {
		t.Error("Expected validation error")
	}

	entryList, err := config.NewServerConfig(config.EntryList, "[CAR_0]\nMODEL=rss_formula_hybrid_2017")
	if err != nil {
		t.Fatal(err)
	}
	cfgs = config.ServerConfigFiles{
		config.ServerCfg: *srvCfg,
		config.EntryList: *entryList,
	}
	_, err = config.NewServerConfigs(cfgs)
	if err != nil {
		t.Error("Unexpected validation error", err)
	}
}

func TestGetServerConfigUnknown(t *testing.T) {
	cm := NewMemoryConfigManager()
	_, err := cm.GetServerConfig("doesnotexist")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestGetServerConfig(t *testing.T) {
	cm := NewMemoryConfigManager()
	_, err := cm.GetServerConfig("doesnotexist")
	if err == nil {
		t.Error("Expected error")
	}
}
