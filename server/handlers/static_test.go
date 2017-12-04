package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/venyii/acfg/server/static"
)

func TestClientAppHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ClientAppHandler(http.FileServer(static.HTTP), static.HTTP))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := []string{
		`<div id="app"></div>`,
		`<script type="text/javascript" src="/static/main.`,
		`<link href="/static/main.`,
	}
	for _, str := range expected {
		if !strings.Contains(rr.Body.String(), str) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	}
}
