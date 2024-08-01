package net

import (
	"backend/db"
	"net/http"
	"slices"
)

var nonAuthenticatedEndpoints = []string{"/users/register", "/users/login", "/debug"}

func AuthenticationMiddleware(next http.Handler, database db.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(nonAuthenticatedEndpoints, r.URL.Path) {
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
