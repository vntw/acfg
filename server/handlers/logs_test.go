package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"

	"acfg/ac"
	"acfg/app"
)

func TestLogsHandlerWithNoData(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := ioutil.TempDir("", "TestLogsHandlerWithNoData")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	sl := ac.MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogsHandler(cfg, sl))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogsHandlerWithData(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := ioutil.TempDir("", "TestLogsHandlerWithData")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	instanceUuid := "9cc27bc7-219c-471d-83bb-7344559f166c"
	tmpInstanceDir := filepath.Join(tmpDir, instanceUuid)
	if err = os.Mkdir(tmpInstanceDir, 0755); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}

	logFile := filepath.Join(tmpInstanceDir, "log1.log")
	if err = ioutil.WriteFile(logFile, []byte("_content_"), 0644); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}
	logFileStat, err := os.Stat(logFile)
	if err != nil {
		t.Fatal("could not get file info:", err)
	}

	sl := ac.MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogsHandler(cfg, sl))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := fmt.Sprintf(`[{"instanceUuid":"9cc27bc7-219c-471d-83bb-7344559f166c","time":%d,"logFiles":[{"name":"log1.log","content":""}]}]`, logFileStat.ModTime().Unix())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogHandlerWithData(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestLogHandlerWithData")
	if err != nil {
		t.Fatal("Could not create tmp dir")
	}
	defer os.RemoveAll(tmpDir)

	instanceUuid := "9cc27bc7-219c-471d-83bb-7344559f166c"
	tmpInstanceDir := filepath.Join(tmpDir, instanceUuid)
	if err = os.Mkdir(tmpInstanceDir, 0755); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}

	logFile := filepath.Join(tmpInstanceDir, "log1.log")
	if err = ioutil.WriteFile(logFile, []byte("_content_"), 0644); err != nil {
		t.Fatal("Could not create dummy log file:", err)
	}
	logFileStat, err := os.Stat(logFile)
	if err != nil {
		t.Fatal("could not get file info:", err)
	}

	sl := ac.MemoryServerLogsManager{}
	cfg := app.Config{
		ServerLogsDir: tmpDir,
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/{instanceUuid}", LogHandler(cfg, sl))

	req, err := http.NewRequest("GET", "/"+instanceUuid, nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"instanceUuid":"9cc27bc7-219c-471d-83bb-7344559f166c","time":%d,"logFiles":[{"name":"log1.log","content":"_content_"}]}`, logFileStat.ModTime().Unix())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
