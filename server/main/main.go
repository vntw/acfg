package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-ini/ini"
	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/venyii/acfg/server/ac"
	"github.com/venyii/acfg/server/ac/server"
	"github.com/venyii/acfg/server/app"
	"github.com/venyii/acfg/server/handlers"
	"github.com/venyii/acfg/server/static"
	"github.com/venyii/acfg/server/user"
)

func main() {
	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	user.SetTokenSecret([]byte(cfg.JwtSecret))
	user.AddConfigUsers(cfg.Users)

	mode := "DEV"
	if cfg.IsProd {
		mode = "PROD"
	}

	ini.PrettyFormat = false

	log.Printf("Running in %s mode on http://0.0.0.0:%d...\n", mode, cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Port), createRouter(cfg)))
}

func createRouter(cfg app.Config) http.Handler {
	im := ac.NewMemoryInstanceManager()
	si := server.InstanceInfoer{}
	sl := ac.MemoryServerLogsManager{}
	cm := ac.NewMemoryConfigManager()

	r := mux.NewRouter()

	r.HandleFunc("/api/sessions", handlers.LoginHandler).Methods("POST")

	r.Handle("/api/logs", auth(handlers.LogsHandler(cfg, sl))).Methods("GET")
	r.Handle("/api/logs/{instanceUuid}", auth(handlers.LogHandler(cfg, sl))).Methods("GET")

	r.Handle("/api/configs", auth(handlers.ServerConfigsHandler(cm))).Methods("GET")
	r.Handle("/api/configs/upload", auth(handlers.ConfigsUploadHandler(cm))).Methods("POST")
	r.Handle("/api/configs/{uuid}/delete", auth(handlers.DeleteConfigHandler(cm))).Methods("DELETE")

	r.Handle("/api/servers", auth(handlers.ServersHandler(im, si, cm))).Methods("GET")
	r.Handle("/api/servers/start", auth(handlers.StartServerHandler(cfg, im, si, cm))).Methods("POST")
	r.Handle("/api/servers/start/upload", auth(handlers.UploadAndStartServerHandler(cfg, im, si, cm))).Methods("POST")
	r.Handle("/api/servers/{uuid}/stop", auth(handlers.StopServerHandler(im))).Methods("POST")
	r.Handle("/api/servers/{uuid}/reconfig", auth(handlers.ReconfigServerHandler(cfg, im, si, cm))).Methods("POST")

	// Client App
	clientAppHandler := handlers.ClientAppHandler(http.FileServer(static.HTTP), static.HTTP)
	r.PathPrefix("/").Handler(http.StripPrefix("/", clientAppHandler))

	var origins []string
	if !cfg.IsProd {
		origins = append(origins, "http://localhost:3000")
	}

	cors := muxHandlers.CORS(
		muxHandlers.AllowedOrigins(origins),
		muxHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "DELETE"}),
		muxHandlers.AllowedHeaders([]string{"Accept", "Authorization", "Accept-Language", "Content-Language", "Origin"}),
		muxHandlers.AllowCredentials(),
	)

	return cors(r)
}

func auth(next http.HandlerFunc) http.Handler {
	return handlers.AuthMiddleware(http.HandlerFunc(next))
}
