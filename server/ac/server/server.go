package server

const (
	StateOffline = iota
	StateOnline
)

type Status struct {
	Ip           string   `json:"ip"`
	Port         int      `json:"port"`
	HttpPort     int      `json:"cport"`
	Name         string   `json:"name"`
	Clients      int      `json:"clients"`
	MaxClients   int      `json:"maxclients"`
	TimeLeft     int      `json:"timeleft"`
	TimeOfDay    int      `json:"timeofday"`
	Pickup       bool     `json:"pickup"`
	Pass         bool     `json:"pass"`
	Session      int      `json:"session"`
	SessionTypes []int    `json:"sessiontypes"`
	Durations    []int    `json:"durations"`
	Track        string   `json:"track"`
	Cars         []string `json:"cars"`
}
