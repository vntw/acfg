package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	IsProd         bool              `envconfig:"ACFG_IS_PROD" default:"true"`
	Port           uint16            `envconfig:"ACFG_PORT" default:"1337"`
	ServerIP       string            `envconfig:"ACFG_SERVER_IP" required:"true"`
	ServerCfgsDir  string            `envconfig:"ACFG_SERVER_CFGS_DIR" default:"server_cfgs"`
	ServerLogsDir  string            `envconfig:"ACFG_SERVER_LOGS_DIR" default:"server_logs"`
	ACServerDir    string            `envconfig:"ACFG_ACSERVER_DIR" required:"true"`
	ACServerBinary string            `envconfig:"ACFG_ACSERVER_BINARY" default:"acServer"`
	StrackerDir    string            `envconfig:"ACFG_STRACKER_DIR"`
	JwtSecret      string            `envconfig:"ACFG_JWT_SECRET" required:"true"`
	Users          map[string]string `envconfig:"ACFG_USERS" required:"true"`
}

func NewConfig() (Config, error) {
	var ec Config

	err := envconfig.Process("", &ec)
	if err != nil {
		return Config{}, errors.New("error processing env config: " + err.Error())
	}

	if len(ec.Users) == 0 {
		return Config{}, errors.New("at least one user must be specified")
	}

	// Validate ServerCfgsDir
	if _, err = os.Stat(ec.ServerCfgsDir); err != nil {
		return Config{}, errors.New("could not find server cfgs dir")
	}
	absServerCfgsDir, err := filepath.Abs(ec.ServerCfgsDir)
	if err != nil {
		return Config{}, errors.New("could not get absolute server cfgs dir")
	}
	ec.ServerCfgsDir = absServerCfgsDir

	// Validate ServerLogsDir
	if _, err = os.Stat(ec.ServerLogsDir); err != nil {
		return Config{}, errors.New("could not find server logs dir")
	}
	absServerLogsDir, err := filepath.Abs(ec.ServerLogsDir)
	if err != nil {
		return Config{}, errors.New("could not get absolute server logs dir")
	}
	ec.ServerLogsDir = absServerLogsDir

	// Validate ACServerDir
	if _, err = os.Stat(ec.ACServerDir); err != nil {
		return Config{}, errors.New("could not find ac server dir")
	}
	absServerPath, err := filepath.Abs(ec.ACServerDir)
	if err != nil {
		return Config{}, errors.New("could not get absolute acServer dir")
	}
	ec.ACServerDir = absServerPath

	// Validate ACServerBinary
	binPath := filepath.Join(ec.ACServerDir, ec.ACServerBinary)
	if _, err = os.Stat(binPath); err != nil {
		return Config{}, errors.New("could not find acServer binary")
	}

	if ec.StrackerDir != "" {
		// Validate StrackerDir
		if _, err = os.Stat(ec.StrackerDir); err != nil {
			return Config{}, errors.New("could not find ac stracker dir")
		}
		absStrackerPath, err := filepath.Abs(ec.StrackerDir)
		if err != nil {
			return Config{}, errors.New("could not get absolute stracker dir")
		}
		ec.StrackerDir = absStrackerPath
	}

	if strings.TrimSpace(ec.JwtSecret) == "" {
		return Config{}, errors.New("JWT secret not set")
	}

	return ec, nil
}
