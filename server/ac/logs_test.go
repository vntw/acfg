package ac

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/venyii/acsrvmanager/server/ac/config"
	instlog "github.com/venyii/acsrvmanager/server/ac/server/log"
	"github.com/venyii/acsrvmanager/server/app"
)

func TestNewServerLog(t *testing.T) {
	sl := instlog.NewLogFile(config.ServerCfg)

	if sl.Name != config.ServerCfg {
		t.Error("Unexpected name")
	}
}

func TestGetServerLogsWithNoData(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestGetServerLogsWithNoData")
	if err != nil {
		t.Fatal("Could not create tmp dir:", err)
	}
	defer os.RemoveAll(tmpDir)

	slm := MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}
	logs, err := slm.GetServerLogs(cfg)

	if len(logs) != 0 {
		t.Fatalf("got wrong server log count: got %v want 0", len(logs))
	}
}

func TestGetServerLogs(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestGetServerLogs")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	instanceUuid := "9cc27bc7-219c-471d-83bb-7344559f166c"
	tmpInstanceDir := filepath.Join(tmpDir, instanceUuid)
	if err = os.Mkdir(tmpInstanceDir, 0755); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}

	if err = ioutil.WriteFile(filepath.Join(tmpInstanceDir, "log1.log"), []byte("_content_"), 0644); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}

	slm := MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}
	logs, err := slm.GetServerLogs(cfg)

	if len(logs) != 1 {
		t.Fatalf("got wrong server log count: got %v want 1", len(logs))
	}

	if logs[0].Files[0].Name != "log1.log" {
		t.Fatalf("got wrong server log name: got %v want log1.log", logs[0].Files[0].Name)
	}

	if logs[0].Files[0].Content != "" {
		t.Fatalf("expected to have no content in list response, got %v", logs[0].Files[0].Content)
	}
}

func TestGetServerLogNotFound(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestGetServerLogNotFound")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	slm := MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}
	_, err = slm.GetServerLog(cfg, "does-not-exist")

	if err == nil {
		t.Fatal("did not expect to find log")
	}

	if !strings.HasSuffix(err.Error(), "/does-not-exist: no such file or directory") {
		t.Fatalf("unexpected error msg, got %v", err.Error())
	}
}

func TestGetServerLog(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestGetServerLog")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	instanceUuid := "9cc27bc7-219c-471d-83bb-7344559f166c"
	tmpInstanceDir := filepath.Join(tmpDir, instanceUuid)
	if err = os.Mkdir(tmpInstanceDir, 0755); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}

	if err = ioutil.WriteFile(filepath.Join(tmpInstanceDir, "log1.log"), []byte("_content_"), 0644); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}

	slm := MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}
	log, err := slm.GetServerLog(cfg, instanceUuid)

	if err != nil {
		t.Fatalf("did not expect to get an error, got %v", err)
	}

	if log.Files[0].Name != "log1.log" {
		t.Fatalf("got wrong server log name: got %v want log1.log", log.Files[0].Name)
	}

	if log.Files[0].Content != "_content_" {
		t.Fatalf("expected to have no content in list response, got %v", log.Files[0].Content)
	}
}
