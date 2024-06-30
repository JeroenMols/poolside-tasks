package net

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		origin string
	}{
		{"*"},
		{"http://localhost:3000"},
	}

	for _, tt := range tests {
		t.Run("Origin "+tt.origin, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), tt.origin).ServeHTTP(w, r)

			assert.Equal(t, http.Header{
				"Access-Control-Allow-Credentials": []string{"true"},
				"Access-Control-Allow-Headers":     []string{"Authorization, Content-Type"},
				"Access-Control-Allow-Methods":     []string{"GET, POST, PUT"},
				"Access-Control-Allow-Origin":      []string{tt.origin},
				"Content-Type":                     []string{"application/json"},
			}, w.Result().Header)
		})
	}
}

func TestCorsMiddleware_Preflight(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/", nil)
	r.Method = http.MethodOptions

	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	CorsMiddleware(fakeHandler, "*").ServeHTTP(w, r)

	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
}
