package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-ini/ini"

	"acfg/ac"
	"acfg/ac/config"
	"acfg/ac/server"
)

func TestServersHandlerWithNoData(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	im := ac.NewMemoryInstanceManager()
	si := server.InstanceInfoer{}
	cm := ac.NewMemoryConfigManager()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServersHandler(im, si, cm))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"configs":[],"servers":[]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestServerConfigsHandlerWithNoData(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cm := ac.NewMemoryConfigManager()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServerConfigsHandler(cm))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestServerConfigsHandlerWithData(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	srvCfg, err := ini.Load([]byte("[SERVER]\nNAME=Test"))
	if err != nil {
		t.Fatal(err)
	}
	elCfg, err := ini.Load([]byte("[CAR_0]\nMODEL=rss_formula_hybrid_2017"))
	if err != nil {
		t.Fatal(err)
	}
	cm := ac.NewMemoryConfigManager()
	cfgs := config.ServerConfigFiles{
		config.ServerCfg: config.ServerConfig{
			Name:     config.ServerCfg,
			Ini:      &config.IniFile{File: srvCfg},
			Checksum: "_checksum1_",
		},
		config.EntryList: config.ServerConfig{
			Name:     config.EntryList,
			Ini:      &config.IniFile{File: elCfg},
			Checksum: "_checksum2_",
		},
	}

	_, err = cm.CreateServerConfigs(cfgs, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(cm.GetServerConfigs()) != 1 {
		t.Fatalf("Expected to have 1 server config, got %v", len(cm.GetServerConfigs()))
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServerConfigsHandler(cm))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var m []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &m)
	if err != nil {
		t.Fatal(err)
	}

	if len(m) != 1 {
		t.Fatalf("handler returned unexpected config count, got %v want 1", len(m))
	}

	// TODO: Complex JSON matching
	expected := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	if !expected.MatchString(m[0]["uuid"].(string)) {
		t.Errorf("handler returned unexpected body: got %v want %v", m[0]["uuid"], expected)
	}
}
