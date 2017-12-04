package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-ini/ini"
	"github.com/gorilla/mux"

	"github.com/venyii/acfg/server/ac"
	"github.com/venyii/acfg/server/ac/config"
)

func TestConfigsUploadHandler(t *testing.T) {
	postData :=
		`--xxx
Content-Disposition: form-data; name="configs"; filename="server_cfg.ini"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

[SERVER]
NAME=Test
--xxx
Content-Disposition: form-data; name="configs"; filename="entry_list.ini"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

[CAR_0]
MODEL=rss_formula_hybrid_2017
--xxx--
`
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	tmpDir, err := ioutil.TempDir("", "TestConfigsUploadHandler")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	cm := ac.NewMemoryConfigManager()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ConfigsUploadHandler(cm))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// TODO: JSON matching
}

func TestConfigsUploadHandlerMissingFile(t *testing.T) {
	postData :=
		`--xxx
Content-Disposition: form-data; name="configs"; filename="server_cfg.ini"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

NAME=Test
--xxx
`
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(strings.NewReader(postData)),
	}

	tmpDir, err := ioutil.TempDir("", "TestConfigsUploadHandler")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	cm := ac.NewMemoryConfigManager()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ConfigsUploadHandler(cm))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"type":"error","message":"entry_list missing"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteConfigHandler(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestConfigsUploadHandler")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	srvCfg, err := ini.Load([]byte("[SERVER]\nNAME=Test"))
	if err != nil {
		t.Fatal(err)
	}
	elCfg, err := ini.Load([]byte("[CAR_0]\nMODEL=rss_formula_hybrid_2017"))
	if err != nil {
		t.Fatal(err)
	}

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

	cm := ac.NewMemoryConfigManager()
	sc, err := cm.CreateServerConfigs(cfgs, true)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/"+sc.Uuid, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/{uuid}", DeleteConfigHandler(cm))
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestDeleteConfigHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/uuid", nil)
	if err != nil {
		t.Fatal(err)
	}

	cm := ac.NewMemoryConfigManager()

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/{uuid}", DeleteConfigHandler(cm))
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"type":"error","message":"Could not find config"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
