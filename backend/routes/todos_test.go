package routes

import (
	"backend/db"
	"backend/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type createTodoTestCase struct {
	description   string
	body          string
	responseCode  int
	responseBody  string
	databaseLists map[string][]models.TodoItem
}

func TestTodo_CreateValidations(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"


	tests := []createListTestCase{
		{
			description:   "Invalid body",
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description:   "Access token not a uuid",
			body:          `{"access_token":"not_a_uuid", "description":"test"}`,
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"invalid access token"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description:   "Access token does not exist",
			body:          fmt.Sprintf(`{"access_token":"%s", "description":"test"}`, nonExistingAccessToken),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description:   "Description too long",
			body:          fmt.Sprintf(`{"access_token":"%s", "description":"%s"}`, validAccessToken, strings.Repeat("a", 257)),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description:   "Description invalid characters",
			body:          fmt.Sprintf(`{"access_token":"%s", "description":"/"}`, validAccessToken),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			database.AccessTokens[validAccessToken] = existingAccount

			todos := Todos{
				Database:     database,
				GenerateUuid: func() string { return "static_uuid" },
			}

			request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(tt.body))
			writer := httptest.NewRecorder()

			todos.Create(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseLists, database.TodoLists)
		})
	}
}
