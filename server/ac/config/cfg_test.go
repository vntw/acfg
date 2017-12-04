package config

import (
	"encoding/json"
	"testing"
)

var simpleServerCfg = `[SERVER]
NAME=Assetto Corsa Server
CARS=rss_formula_hybrid_2017;rss_formula_hybrid_2017_s1
CONFIG_TRACK=
TRACK=spa
SUN_ANGLE=16
UDP_PORT=9456
TCP_PORT=9457
HTTP_PORT=8098
MAX_BALLAST_KG=0
QUALIFY_MAX_WAIT_PERC=120
RACE_PIT_WINDOW_START=0
RACE_PIT_WINDOW_END=0
[DATA]
DESCRIPTION=Hello World!
`

func TestNewServerConfig(t *testing.T) {
	sc, err := NewServerConfig(ServerCfg, simpleServerCfg)
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if sc.Name != ServerCfg {
		t.Error("Unexpected name")
	}
	if sc.Checksum == "" {
		t.Error("Checksum should not be empty")
	}
	if sc.Ini == nil {
		t.Error("Ini should not be null")
	}
}

func TestMarshalServerConfig(t *testing.T) {
	sc, err := NewServerConfig(ServerCfg, simpleServerCfg)
	if err != nil {
		t.Fatal("unexpected error")
	}

	j, err := json.Marshal(sc)
	if err != nil {
		t.Fatal("unexpected error while marshaling server config", err)
	}

	expected := `{"name":"server_cfg","ini":{"SERVER":{"NAME":"Assetto Corsa Server","CARS":"rss_formula_hybrid_2017;rss_formula_hybrid_2017_s1","CONFIG_TRACK":"","TRACK":"spa","SUN_ANGLE":"16","UDP_PORT":"9456","TCP_PORT":"9457","HTTP_PORT":"8098","MAX_BALLAST_KG":"0","QUALIFY_MAX_WAIT_PERC":"120","RACE_PIT_WINDOW_START":"0","RACE_PIT_WINDOW_END":"0"},"DATA":{"DESCRIPTION":"Hello World!"}},"checksum":"6f59319db8356ca6dcecd1b671007371"}`
	if string(j) != expected {
		t.Fatalf("unexpected marshaled server config, got\n%v want\n%v", string(j), expected)
	}
}
