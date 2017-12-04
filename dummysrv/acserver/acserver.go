package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-ini/ini"
)

// server_cfg.ini
type ServerCfg struct {
	Server serverCfgServer `ini:"SERVER"`
}

// server_cfg.ini [SERVER] section
type serverCfgServer struct {
	Name     string `ini:"NAME"`
	HttpPort int    `ini:"HTTP_PORT"`
}

func main() {
	serverCfg := flag.String("c", "cfg/server_cfg.ini", "Location for cfg file")
	entryList := flag.String("e", "cfg/entry_list.ini", "Location for entry list file")

	flag.Parse()

	assertCfgExists(*serverCfg)
	assertCfgExists(*entryList)

	x := new(ServerCfg)
	err := ini.MapTo(x, *serverCfg)
	if err != nil {
		log.Fatalln("could not map ini from server_cfg.ini")
	}

	log.Printf("Using server_cfg: %s\n", *serverCfg)
	log.Printf("Using entry_list: %s\n", *entryList)

	http.HandleFunc("/INFO", infoHandler(x))
	http.HandleFunc("/JSON|", jsonHandler)
	http.HandleFunc("/ENTRY", entryHandler)
	http.HandleFunc("/", homeHandler)

	log.Printf("Running on http://localhost:%d\n", x.Server.HttpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", x.Server.HttpPort), nil); err != nil {
		log.Fatalf("acServer: %v", err)
	}
}

func infoHandler(srvCfg *ServerCfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json, err := ioutil.ReadFile("responses/INFO.json")
		if err != nil {
			log.Fatalln("Could not find INFO.json file")
		}

		tmpl, err := template.New("info").Parse(string(json))
		if err != nil {
			log.Fatalln("Could not parse INFO.json template file")
		}

		w.Header().Set("Content-Type", "application/json")
		err = tmpl.Execute(w, srvCfg)
		if err != nil {
			log.Fatalln("could not exec template file", err)
		}
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	json, err := ioutil.ReadFile("responses/JSON|.json")
	if err != nil {
		log.Fatalln("Could not find JSON|.json file")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("responses/ENTRY.html")
	if err != nil {
		log.Fatalln("Could not find ENTRY.html file")
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
  <title>AC Dummy Server</title>
</head>
<body>
  <h1>AC Dummy Server</h1>
  <a href="/INFO">/INFO</a>
  <a href="/JSON|">/JSON|</a>
  <a href="/ENTRY">/ENTRY</a>
</body>
</html>
`))
}

func assertCfgExists(cfg string) {
	if _, err := os.Stat(cfg); err != nil {
		log.Fatalln(err)
	}
}
