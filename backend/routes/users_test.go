package routes

import (
	"backend/db"
	"backend/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type registerTestCase struct {
	description   string
	body          string
	responseCode  int
	responseBody  string
	databaseUsers map[string]db.User
}

func TestUsers_Register(t *testing.T) {

	tests := []registerTestCase{
		{
			description:   "Valid body",
			body:          `{"name":"myname"}`,
			responseCode:  http.StatusOK,
			responseBody:  `{"account_number":"static_uuid"}`,
			databaseUsers: map[string]db.User{"static_uuid": {AccountNumber: "static_uuid", Name: "myname"}},
		},
		{
			description:   "Invalid body",
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseUsers: make(map[string]db.User),
		},
		{
			description:   "User name too short",
			body:          `{"name":"s"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid user name"}`,
			databaseUsers: make(map[string]db.User),
		},
		{
			description:   "User name invalid character",
			body:          `{"name":"name-%*("}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid user name"}`,
			databaseUsers: make(map[string]db.User),
		},
		{
			description:   "User name too long",
			body:          fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("a", 33)),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid user name"}`,
			databaseUsers: make(map[string]db.User),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func() string { return "static_uuid" },
			)
			users := CreateUsers(database)

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
	databaseTokens map[string]db.AccessToken
}

func TestUsers_Login(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-e0d154b5d433"
	const nonExistingAccount = "f2d869a8-e5bc-4fbf-ad71-e0d154b5d434"

	tests := []loginTestCase{
		{
			description:    "Valid body",
			body:           fmt.Sprintf(`{"account_number":"%s"}`, existingAccount),
			responseCode:   http.StatusOK,
			responseBody:   `{"access_token":"static_uuid"}`,
			databaseTokens: map[string]db.AccessToken{"static_uuid": {AccountNumber: existingAccount, Token: "static_uuid"}},
		},
		{
			description:    "Invalid body",
			body:           `{"invalid":"body"}`,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"error":"invalid body"}`,
			databaseTokens: make(map[string]db.AccessToken),
		},
		{
			description:    "Account number not a uuid",
			body:           `{"account_number":"not_a_uuid"}`,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"error":"invalid account number"}`,
			databaseTokens: make(map[string]db.AccessToken),
		},
		{
			description:    "Account does not exist",
			body:           fmt.Sprintf(`{"account_number":"%s"}`, nonExistingAccount),
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"error":"account not found"}`,
			databaseTokens: make(map[string]db.AccessToken),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2021, 1, 1) },
				func() string { return "static_uuid" },
			)
			database.Users[existingAccount] = db.User{AccountNumber: existingAccount, Name: "myname"}
			users := CreateUsers(database)

			request := httptest.NewRequest(http.MethodGet, "/users/login", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			users.Login(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseTokens, database.AccessTokens)
		})
	}
}
