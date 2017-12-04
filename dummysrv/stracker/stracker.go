package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/go-ini/ini"
)

type strackerCfg struct {
	S serverSection `ini:"STRACKER_CONFIG"`
	H httpSection   `ini:"HTTP_CONFIG"`
}

// [STRACKER_CONFIG] ac_server_cfg_ini listening_port
type serverSection struct {
	IniPath string `ini:"ac_server_cfg_ini"`
	Port    int    `ini:"listening_port"`
}

// [HTTP_CONFIG] listen_port
type httpSection struct {
	Port int `ini:"listen_port"`
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("usage: stracker --stracker_ini path/to/stracker.ini")
	}

	strackerIniPath := os.Args[2]
	if _, err := os.Stat(strackerIniPath); err != nil {
		log.Fatalln("could not find stracker.ini", err)
	}

	log.Printf("Using stracker.ini: %s\n", strackerIniPath)

	x := new(strackerCfg)
	err := ini.MapTo(x, strackerIniPath)
	if err != nil {
		log.Fatalln("invalid stracker.ini format:", err)
	}

	log.Printf("AC Cfg:\t%s\n", x.S.IniPath)
	log.Printf("TCP Port:\t%d\n", x.S.Port)
	log.Printf("HTTP Port:\t%d\n", x.H.Port)

	_, err = net.Listen("tcp", ":"+strconv.Itoa(x.S.Port))
	if err != nil {
		log.Fatalln("could not listen:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ACFG Dummy Stracker"))
	})

	log.Printf("Running on http://localhost:%d\n", x.H.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", x.H.Port), nil)
}
