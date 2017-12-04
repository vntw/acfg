package config

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-ini/ini"
	"github.com/google/uuid"
)

const (
	ServerCfg = "server_cfg"
	EntryList = "entry_list"
)

type IniFile struct {
	*ini.File
}

// Preserve ini ordering when marshaling to json
func (f IniFile) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	sections := f.Sections()
	secLength := len(sections)
	secCount := 0
	for _, sec := range sections {
		if sec.Name() == ini.DEFAULT_SECTION {
			secCount++
			continue
		}

		buffer.WriteString(fmt.Sprintf("\"%s\":{", sec.Name()))

		keys := sec.Keys()
		keyLength := len(keys)
		keyCount := 0
		for _, key := range keys {
			jsonValue, err := json.Marshal(key.Value())
			if err != nil {
				return nil, err
			}
			buffer.WriteString(fmt.Sprintf("\"%s\":%s", key.Name(), string(jsonValue)))
			keyCount++
			if keyCount < keyLength {
				buffer.WriteString(",")
			}
		}

		buffer.WriteString("}")

		secCount++
		if secCount < secLength {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

type ServerConfig struct {
	Name     string   `json:"name"`
	Ini      *IniFile `json:"ini"`
	Checksum string   `json:"checksum"`
}

func (sc ServerConfig) Filename() string {
	return fmt.Sprintf("%s.ini", sc.Name)
}

func NewServerConfig(name string, content string) (*ServerConfig, error) {
	sc := &ServerConfig{Name: name}
	sc.Checksum = fmt.Sprintf("%x", md5.Sum([]byte(content)))

	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, []byte(content))
	if err != nil {
		return nil, err
	}

	sc.Ini = &IniFile{cfg}

	return sc, nil
}

type ServerConfigFiles map[string]ServerConfig

func (scf ServerConfigFiles) Validate() error {
	if _, ok := scf[ServerCfg]; !ok {
		return errors.New(ServerCfg + " missing")
	}
	if _, ok := scf[EntryList]; !ok {
		return errors.New(EntryList + " missing")
	}
	if scf[ServerCfg].Checksum == "" {
		return errors.New(ServerCfg + " checksum missing")
	}
	if scf[EntryList].Checksum == "" {
		return errors.New(EntryList + " checksum missing")
	}

	ssec, err := scf[ServerCfg].Ini.GetSection("SERVER")
	if err != nil || !ssec.HasKey("NAME") {
		return errors.New("invalid " + ServerCfg + " config content")
	}

	esec, err := scf[EntryList].Ini.GetSection("CAR_0")
	if err != nil || !esec.HasKey("MODEL") {
		return errors.New("invalid " + EntryList + " config content")
	}

	return nil
}

type ServerConfigs struct {
	Uuid  string            `json:"uuid"`
	Name  string            `json:"name"`
	Files ServerConfigFiles `json:"files"`
}

func (sc ServerConfigs) Path(dir, cfg string) string {
	return fmt.Sprintf("%s/%s/%s.ini", dir, sc.Uuid, cfg)
}

func NewServerConfigs(cfgs ServerConfigFiles) (ServerConfigs, error) {
	if err := cfgs.Validate(); err != nil {
		return ServerConfigs{}, err
	}

	return ServerConfigs{
		Uuid:  uuid.New().String(),
		Name:  cfgs[ServerCfg].Ini.Section("SERVER").Key("NAME").Value(),
		Files: cfgs,
	}, nil
}
