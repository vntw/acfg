package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"acfg/ac"
	"acfg/ac/plugins"
	"acfg/ac/server"
	"acfg/ac/spec"
	"acfg/app"
)

func ServersHandler(im ac.InstanceManager, si server.InstanceInfoer, cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		servers := []server.ServerInstance{}
		for _, inst := range im.GetInstances() {
			servers = append(servers, si.GetServerInstance(inst))
		}

		data := map[string]interface{}{
			"servers": servers,
			"configs": cm.GetServerConfigs(),
		}
		sendResponse(w, data)
	}
}

func ServerConfigsHandler(cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendResponse(w, cm.GetServerConfigs())
	}
}

func StartServerHandler(cfg app.Config, im ac.InstanceManager, si server.InstanceInfoer, cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverCfgUuid := r.PostFormValue("serverCfgUuid")

		sc, err := cm.GetServerConfig(serverCfgUuid)
		if err != nil {
			sendError(w, err, fmt.Sprintf("Could not find server config '%s'", serverCfgUuid), 0)
			return
		}

		srvSpec := spec.NewServerSpec()
		srvSpec.Preset = sc
		srvSpec.AddPlugin(plugins.NewStrackerPlugin(cfg.StrackerDir))

		instance, err := im.StartInstanceWithConfig(cfg, srvSpec)
		if err != nil {
			sendError(w, err, "Could not start instance", 0)
			return
		}

		sendResponse(w, si.GetServerInstance(instance))
	}
}

func StopServerHandler(im ac.InstanceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]

		if err := im.StopInstance(uuid); err != nil {
			sendError(w, err, "Could not stop instance", 0)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ReconfigServerHandler(cfg app.Config, im ac.InstanceManager, si server.InstanceInfoer, cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]

		instance, err := im.GetInstance(uuid)
		if err != nil {
			sendError(w, err, "Could not find instance", 0)
			return
		}

		cfgs, err := handleConfigFilesUpload(r)
		if err != nil {
			sendError(w, err, err.Error(), 0)
			return
		}

		if err := im.StopInstance(uuid); err != nil {
			sendError(w, err, "Could not stop instance", 0)
			return
		}

		sc, err := cm.CreateServerConfigs(cfgs, false)
		if err != nil {
			sendError(w, err, err.Error(), 400)
			return
		}

		newSpec := &instance.Spec
		newSpec.Preset = sc

		tmpsc, err := ac.CreateTmpConfig(cfg.ServerCfgsDir, newSpec)
		if err != nil {
			sendError(w, err, err.Error(), 400)
			return
		}

		newSpec.TmpConfig = tmpsc

		instance, err = im.StartInstance(cfg, newSpec)
		if err != nil {
			sendError(w, err, "Could not start instance", 0)
			return
		}

		sendResponse(w, si.GetServerInstance(instance))
	}
}

func UploadAndStartServerHandler(cfg app.Config, im ac.InstanceManager, si server.InstanceInfoer, cm ac.ConfigManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfgs, err := handleConfigFilesUpload(r)
		if err != nil {
			sendError(w, err, err.Error(), 0)
			return
		}

		sc, err := cm.CreateServerConfigs(cfgs, false)
		if err != nil {
			sendError(w, err, err.Error(), 400)
			return
		}

		srvSpec := spec.NewServerSpec()
		srvSpec.Preset = sc
		srvSpec.AddPlugin(plugins.NewStrackerPlugin(cfg.StrackerDir))

		instance, err := im.StartInstanceWithConfig(cfg, srvSpec)
		if err != nil {
			sendError(w, err, "Could not start instance", 0)
			return
		}

		sendResponse(w, si.GetServerInstance(instance))
	}
}
