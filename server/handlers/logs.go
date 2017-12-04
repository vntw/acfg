package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/venyii/acsrvmanager/server/ac"
	"github.com/venyii/acsrvmanager/server/app"
)

func LogsHandler(cfg app.Config, sl ac.ServerLogsManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs, err := sl.GetServerLogs(cfg)

		if err != nil {
			sendError(w, err, "Could not get logs", 0)
			return
		}

		sendResponse(w, logs)
	}
}

func LogHandler(cfg app.Config, sl ac.ServerLogsManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceUuid := mux.Vars(r)["instanceUuid"]

		logs, err := sl.GetServerLog(cfg, instanceUuid)

		if err != nil {
			sendError(w, err, "Could not get log", 0)
			return
		}

		sendResponse(w, logs)
	}
}
