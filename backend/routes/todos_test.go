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
	"time"
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
	const existingList = "f2d869a8-e5bc-4fbf-ad71-444444444444"

	tests := []createListTestCase{
		{
			description:   "Invalid body",
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description: "Access token not a uuid",
			body: fmt.Sprintf(
				`{"access_token":"%s", "description":"%s", "todo_list_id": "%s"}`,
				"invalid_token", "fake_description", existingList),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"invalid access token"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description: "Access token does not exist",
			body: fmt.Sprintf(
				`{"access_token":"%s", "description":"%s", "todo_list_id": "%s"}`,
				nonExistingAccessToken, "fake_description", existingList),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description: "Description too long",
			body: fmt.Sprintf(`{"access_token":"%s", "description":"%s", "todo_list_id":"%s"}`,
				validAccessToken, strings.Repeat("a", 257), existingList),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseLists: make(map[string][]models.TodoItem),
		},
		{
			description: "Description invalid characters",
			body: fmt.Sprintf(`{"access_token":"%s", "description":"%s", "todo_list_id":"%s"}`,
				validAccessToken, "/", existingList),
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
				CurrentTime:  func() time.Time { return fixedTestTime() },
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

func TestTodo_CreateLogic(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"
	const existingList = "f2d869a8-e5bc-4fbf-ad71-444444444444"
	const nonExistingList = "f2d869a8-e5bc-4fbf-ad71-555555555555"

	tests := []createListTestCase{
		{
			description: "Create new todo",
			body: fmt.Sprintf(`{"access_token":"%s", "description":"%s", "todo_list_id":"%s"}`,
				validAccessToken, "test todo", existingList),
			responseCode: http.StatusOK,
			responseBody: `{"created_by":"","description":"test todo","status":"todo","updated_at":"2024-06-30T00:00:00+00:00"}`,
			databaseLists: map[string][]models.TodoItem{
				existingList: {models.TodoItem{
					Description: "test todo",
					Status:      "todo",
					User:        existingAccount,
					UpdatedAt:   fixedTestTime(),
				}},
			},
		},
		{
			description: "Todo list does not exist",
			body: fmt.Sprintf(`{"access_token":"%s", "description":"%s", "todo_list_id":"%s"}`,
				validAccessToken, "test todo", nonExistingList),
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"Todo list not found"}`,
			databaseLists: map[string][]models.TodoItem{
				existingList: {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[existingList] = []models.TodoItem{}

			todos := Todos{
				Database:     database,
				GenerateUuid: func() string { return "static_uuid" },
				CurrentTime:  func() time.Time { return fixedTestTime() },
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

func fixedTestTime() time.Time {
	return time.Date(2024, 6, 30, 0, 0, 0, 0, time.FixedZone("CEST", 1))
}
