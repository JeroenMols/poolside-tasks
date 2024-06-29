package routes

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testCase struct {
	description  string
	body         string
	responseCode int
	responseBody string
}

func TestUsers_Register(t *testing.T) {

	tests := []testCase{
		{
			description:  "Missing body",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: "{\"error\":\"Invalid body\"}",
		},
		{
			description:  "Valid body",
			body:         "{\"name\":\"myname\"}",
			responseCode: http.StatusOK,
			responseBody: "{\"account_number\":\"static_uuid\"}",
		},
		{
			description:  "Additional body attribute",
			body:         "{\"name\":\"myname\",\"age\":30}",
			responseCode: http.StatusBadRequest,
			responseBody: "{\"error\":\"Invalid body\"}",
		},
		{
			description:  "Invalid body",
			body:         "{\"invalid\":\"body\"}",
			responseCode: http.StatusBadRequest,
			responseBody: "{\"error\":\"Invalid body\"}",
		},
	}

	// TODO enforce min size for user name

	users := Users{
		AddResponseHeaders: func(w http.ResponseWriter) {},
		GenerateUuid: func() string {
			return "static_uuid"
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/users/register", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			users.Register(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
		})
	}
}
