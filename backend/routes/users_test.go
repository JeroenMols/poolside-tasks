package routes

import (
	"backend/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type registerTestCase struct {
	description   string
	body          string
	responseCode  int
	responseBody  string
	databaseUsers map[string]string
}

func TestUsers_Register(t *testing.T) {

	tests := []registerTestCase{
		{
			description:   "Missing body",
			body:          "",
			responseCode:  http.StatusBadRequest,
			responseBody:  "{\"error\":\"invalid body\"}",
			databaseUsers: make(map[string]string),
		},
		{
			description:   "Valid body",
			body:          "{\"name\":\"myname\"}",
			responseCode:  http.StatusOK,
			responseBody:  "{\"account_number\":\"static_uuid\"}",
			databaseUsers: map[string]string{"static_uuid": "myname"},
		},
		{
			description:   "Additional body attribute",
			body:          "{\"name\":\"myname\",\"age\":30}",
			responseCode:  http.StatusBadRequest,
			responseBody:  "{\"error\":\"invalid body\"}",
			databaseUsers: make(map[string]string),
		},
		{
			description:   "Invalid body",
			body:          "{\"invalid\":\"body\"}",
			responseCode:  http.StatusBadRequest,
			responseBody:  "{\"error\":\"invalid body\"}",
			databaseUsers: make(map[string]string),
		},
		{
			description:   "Empty body",
			body:          "{}",
			responseCode:  http.StatusBadRequest,
			responseBody:  "{\"error\":\"validation error\"}",
			databaseUsers: make(map[string]string),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			users := Users{
				Database: database,
				GenerateUuid: func() string {
					return "static_uuid"
				},
			}

			request := httptest.NewRequest(http.MethodGet, "/users/register", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			users.Register(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseUsers, database.Users)
		})
	}
}

type loginTestCase struct {
	description    string
	body           string
	responseCode   int
	responseBody   string
	databaseTokens map[string]string
}

func TestUsers_Login(t *testing.T) {

	tests := []loginTestCase{
		{
			description:    "Missing body",
			body:           "",
			responseCode:   http.StatusBadRequest,
			responseBody:   "{\"error\":\"invalid body\"}",
			databaseTokens: make(map[string]string),
		},
		{
			description:    "Valid body",
			body:           "{\"account_number\":\"my_number\"}",
			responseCode:   http.StatusOK,
			responseBody:   "{\"access_token\":\"static_uuid\"}",
			databaseTokens: map[string]string{"my_number": "static_uuid"},
		},
		{
			description:    "Additional body attribute",
			body:           "{\"account_number\":\"my_number\",\"age\":30}",
			responseCode:   http.StatusBadRequest,
			responseBody:   "{\"error\":\"invalid body\"}",
			databaseTokens: make(map[string]string),
		},
		{
			description:    "Invalid body",
			body:           "{\"invalid\":\"body\"}",
			responseCode:   http.StatusBadRequest,
			responseBody:   "{\"error\":\"invalid body\"}",
			databaseTokens: make(map[string]string),
		},
		{
			description:    "Empty body",
			body:           "{}",
			responseCode:   http.StatusBadRequest,
			responseBody:   "{\"error\":\"validation error\"}",
			databaseTokens: make(map[string]string),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			users := Users{
				Database: database,
				GenerateUuid: func() string {
					return "static_uuid"
				},
			}

			request := httptest.NewRequest(http.MethodGet, "/users/login", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			users.Login(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseTokens, database.AccessTokens)
		})
	}
}
