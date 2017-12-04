package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/venyii/acfg/server/ac"
)

func ConfigsUploadHandler(cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfgs, err := handleConfigFilesUpload(r)
		if err != nil {
			sendError(w, err, err.Error(), 0)
			return
		}

		sc, err := cm.CreateServerConfigs(cfgs, true)
		if err != nil {
			sendError(w, err, err.Error(), 400)
			return
		}

		sendResponse(w, sc)
	}
}

func DeleteConfigHandler(cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]

		if err := cm.DeleteServerConfig(uuid); err != nil {
			sendError(w, err, "Could not find config", 404)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
