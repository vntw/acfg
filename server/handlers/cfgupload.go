package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"acfg/ac/config"
)

func handleConfigFilesUpload(r *http.Request) (config.ServerConfigFiles, error) {
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(50000)
		if err != nil {
			return nil, err
		}
	}

	cfgs := config.ServerConfigFiles{}

	fhs := r.MultipartForm.File["configs"]
	for _, fh := range fhs {
		err := func() error {
			f, err := fh.Open()
			if err != nil {
				return err
			}
			defer f.Close()

			name := strings.Split(filepath.Base(fh.Filename), ".")[0]

			if name != config.ServerCfg && name != config.EntryList {
				return errors.New(fmt.Sprintf("invalid file '%s' uploaded", name))
			}

			var buf bytes.Buffer
			io.Copy(&buf, f)

			cfg, err := config.NewServerConfig(name, buf.String())
			if err != nil {
				return err
			}

			cfgs[name] = *cfg

			return nil
		}()

		if err != nil {
			return nil, err
		}
	}

	return cfgs, nil
}
