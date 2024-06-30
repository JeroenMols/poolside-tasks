package routes

import (
	"backend/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type createListTestCase struct {
	description   string
	body          string
	responseCode  int
	responseBody  string
	databaseLists map[string][]db.TodoItem
}

func TestTodoLists_Create(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"

	tests := []createListTestCase{
		{
			description:   "Invalid body",
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseLists: make(map[string][]db.TodoItem),
		},
		{
			description:   "Access token not a uuid",
			body:          `{"access_token":"not_a_uuid"}`,
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"invalid access token"}`,
			databaseLists: make(map[string][]db.TodoItem),
		},
		{
			description:   "Access token does not exist",
			body:          fmt.Sprintf(`{"access_token":"%s"}`, nonExistingAccessToken),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseLists: make(map[string][]db.TodoItem),
		},
		{
			description:   "Access token does not exist2",
			body:          fmt.Sprintf(`{"access_token":"%s"}`, nonExistingAccessToken),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseLists: make(map[string][]db.TodoItem),
		},
		{
			description:   "Create new todo list",
			body:          fmt.Sprintf(`{"access_token":"%s"}`, validAccessToken),
			responseCode:  http.StatusOK,
			responseBody:  `{"todo_list_id":"static_uuid"}`,
			databaseLists: map[string][]db.TodoItem{"static_uuid": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			database.AccessTokens[validAccessToken] = existingAccount

			todoList := TodoLists{
				Database:     database,
				GenerateUuid: func() string { return "static_uuid" },
			}

			request := httptest.NewRequest(http.MethodPost, "/todolists", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			todoList.Create(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseLists, database.TodoLists)
		})
	}
}
