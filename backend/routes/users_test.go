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
			responseBody:  `{"user_id":"static_uuid"}`,
			databaseUsers: map[string]db.User{"static_uuid": {Id: "static_uuid", Name: "myname"}},
		},
		{
			description:   "Invalid body",
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"body not valid"}`,
			databaseUsers: make(map[string]db.User),
		},
		{
			description:   "UserId name too short",
			body:          `{"name":"s"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid user name"}`,
			databaseUsers: make(map[string]db.User),
		},
		{
			description:   "UserId name invalid character",
			body:          `{"name":"name-%*("}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid user name"}`,
			databaseUsers: make(map[string]db.User),
		},
		{
			description:   "UserId name too long",
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
				func(string) string { return "static_uuid" },
			)
			users := CreateUsers(&database)

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
	tests := []loginTestCase{
		{
			description:    "Valid body",
			body:           fmt.Sprintf(`{"user_id":"%s"}`, fakeUserId),
			responseCode:   http.StatusOK,
			responseBody:   `{"access_token":"static_uuid"}`,
			databaseTokens: map[string]db.AccessToken{"static_uuid": {UserId: fakeUserId, Token: "static_uuid"}},
		},
		{
			description:    "Invalid body",
			body:           `{"invalid":"body"}`,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"error":"body not valid"}`,
			databaseTokens: make(map[string]db.AccessToken),
		},
		{
			description:    "Account number not a uuid",
			body:           `{"user_id":"not_a_uuid"}`,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"error":"invalid user id"}`,
			databaseTokens: make(map[string]db.AccessToken),
		},
		{
			description:    "Account does not exist",
			body:           fmt.Sprintf(`{"user_id":"%s"}`, fakeWrongUserId),
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"error":"user not found"}`,
			databaseTokens: make(map[string]db.AccessToken),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2021, 1, 1) },
				func(string) string { return "static_uuid" },
			)
			database.Users[fakeUserId] = db.User{Id: fakeUserId, Name: "myname"}
			users := CreateUsers(&database)

			request := httptest.NewRequest(http.MethodGet, "/users/login", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			users.Login(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseTokens, database.AccessTokens)
		})
	}
}
