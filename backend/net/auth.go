package net

import (
	"backend/db"
	"net/http"
	"slices"
)

var noAuthenticatedEndpoints = []string{"/users/register", "/users/login"}

func AuthenticationMiddleware(next http.Handler, database db.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(noAuthenticatedEndpoints, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if _, err := database.GetAccessToken(r.Header.Get("Authorization")); err != nil {
			HaltUnauthorized(w, err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}
