package handlers

import (
	"net/http"
	"strings"

	"github.com/venyii/acsrvmanager/server/user"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	u, err := user.MatchUser(username, password)
	if err != nil {
		sendError(w, err, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := user.CreateToken(u)
	if err != nil {
		sendError(w, err, "Could not create token", 500)
		return
	}

	sendResponse(w, map[string]string{"token": token})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" || !user.AuthToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
