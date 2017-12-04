package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func fetchServerInfo(i Instance) (*http.Response, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	return client.Get(fmt.Sprintf("http://localhost:%d/INFO", i.Spec.Ports.HttpPort))
}

type InstanceInfoer struct{}

func (si InstanceInfoer) GetServerInstance(i Instance) ServerInstance {
	resp, err := fetchServerInfo(i)
	if err != nil {
		log.Print("Error getting server status:", err)
		return ServerInstance{
			Error:    "Error getting server status",
			State:    StateOffline,
			Instance: i,
		}
	}
	defer resp.Body.Close()

	status := &Status{}
	if err = json.NewDecoder(resp.Body).Decode(status); err != nil {
		log.Print("Error decoding server status:", err)
		return ServerInstance{
			Error:    "Error decoding server status",
			State:    StateOffline,
			Instance: i,
		}
	}

	return ServerInstance{
		State:    StateOnline,
		Instance: i,
		Status:   *status,
	}
}
