package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewValidConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "acfg_env_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = ioutil.WriteFile(filepath.Join(tmpDir, "acServer"), []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}

	os.Clearenv()

	os.Setenv("ACFG_SERVER_IP", "localhost")
	os.Setenv("ACFG_PORT", "7331")
	os.Setenv("ACFG_SERVER_CFGS_DIR", tmpDir)
	os.Setenv("ACFG_SERVER_LOGS_DIR", tmpDir)
	os.Setenv("ACFG_ACSERVER_DIR", tmpDir)
	os.Setenv("ACFG_USERS", "u1:p1,u2:p2,u3:p3")
	os.Setenv("ACFG_JWT_SECRET", "_secret_")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error for valid env config, got %v", err)
	}

	if len(cfg.Users) != 3 {
		t.Errorf("Expected three configured users, got %v", len(cfg.Users))
	}
}
