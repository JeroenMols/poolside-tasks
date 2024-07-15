package net

import (
	"backend/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationMiddlewareMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		token          string
		expectedStatus int
	}{
		{
			name:           "No auth needed for /user/register",
			url:            "http://localhost:3000/users/register",
			token:          "",
			expectedStatus: http.StatusOK},
		{
			name:           "No auth needed for /user/login",
			url:            "http://localhost:3000/users/login",
			token:          "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized when missing token",
			url:            "http://localhost:3000/todolists/",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Unauthorized when invalid token",
			url:            "http://localhost:3000/todolists/",
			token:          "invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Unauthorized when non existing token",
			url:            "http://localhost:3000/todolists/",
			token:          "tkn_ffffffffffffffffffffff",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Authorized when existing token",
			url:            "http://localhost:3000/todolists/",
			token:          "tkn_aaaaaaaaaaaaaaaaaaaaaa",
			expectedStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run("Origin "+tt.url, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			r.Header.Set("Authorization", tt.token)

			database := db.TestDatabase(nil, nil)
			database.AccessTokens["tkn_aaaaaaaaaaaaaaaaaaaaaa"] = db.AccessToken{UserId: "valid_user_id", Token: "tkn_aaaaaaaaaaaaaaaaaaaaaa"}

			AuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), &database).ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Result().StatusCode)
		})
	}
}
