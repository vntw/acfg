package ac

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	instlog "github.com/venyii/acfg/server/ac/server/log"
	"github.com/venyii/acfg/server/app"
)

type ServerLogsManager interface {
	GetServerLogs(cfg app.Config) ([]instlog.Container, error)
	GetServerLog(cfg app.Config, instanceUuid string) (instlog.Container, error)
}

type MemoryServerLogsManager struct{}

func (msl MemoryServerLogsManager) GetServerLogs(cfg app.Config) ([]instlog.Container, error) {
	logs := make([]instlog.Container, 0)

	files, err := ioutil.ReadDir(cfg.ServerLogsDir)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].ModTime().After(files[j].ModTime()) })

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		files, err := ioutil.ReadDir(filepath.Join(cfg.ServerLogsDir, file.Name()))
		if err != nil {
			return nil, err
		}

		if len(files) <= 0 {
			continue
		}

		l := instlog.NewLogContainer(file.Name(), file.ModTime().Unix())

		for _, file := range files {
			l.Files = append(l.Files, instlog.File{Name: file.Name()})
		}

		logs = append(logs, *l)
	}

	return logs, nil
}

func (msl MemoryServerLogsManager) GetServerLog(cfg app.Config, instanceUuid string) (instlog.Container, error) {
	if instanceUuid == "" {
		return instlog.Container{}, errors.New("instanceUuid must not be empty")
	}

	path, err := filepath.Abs(filepath.Join(cfg.ServerLogsDir, instanceUuid))
	if err != nil {
		return instlog.Container{}, err
	}

	if !strings.HasPrefix(path, cfg.ServerLogsDir) {
		return instlog.Container{}, errors.New("invalid log path given")
	}

	stat, err := os.Stat(path)
	if err != nil {
		return instlog.Container{}, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return instlog.Container{}, err
	}

	l := instlog.NewLogContainer(instanceUuid, stat.ModTime().Unix())

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return instlog.Container{}, err
		}

		l.Files = append(l.Files, instlog.File{Name: file.Name(), Content: string(content)})
	}

	return *l, nil
}
