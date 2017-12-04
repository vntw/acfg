package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/venyii/acfg/server/user"
)

func TestLoginHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	user.AddConfigUsers(map[string]string{
		"u1": "p1",
	})

	req.PostForm = url.Values{}
	req.PostForm.Add("username", "u1")
	req.PostForm.Add("password", "p1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// TODO: Variable value match
	//expected := `{"token":"#token#"}`
	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}
}

func TestLoginHandlerWithWrongLoginCreds(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	user.AddConfigUsers(map[string]string{
		"u1": "p1",
	})

	req.PostForm = url.Values{}
	req.PostForm.Add("username", "u2")
	req.PostForm.Add("password", "p2")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := `{"type":"error","message":"Invalid credentials"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
