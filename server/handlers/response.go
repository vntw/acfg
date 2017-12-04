package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func sendResponse(w http.ResponseWriter, resp interface{}) {
	response, err := json.Marshal(resp)

	if err != nil {
		sendMarshalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func sendError(w http.ResponseWriter, err error, msg string, statusCode int) {
	log.Println(msg, err)

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	resp := Response{"error", msg, statusCode}
	response, err := json.Marshal(resp)

	if err != nil {
		sendMarshalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(response)
}

func sendMarshalError(w http.ResponseWriter, err error) {
	log.Println("Could not marshal response:", err)

	resp := Response{"error", "", http.StatusInternalServerError}

	w.WriteHeader(resp.StatusCode)

	response, err := json.Marshal(resp)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
