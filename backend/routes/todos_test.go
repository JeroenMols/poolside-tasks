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
	accessToken   string
	body          string
	responseCode  int
	responseBody  string
	databaseLists map[string][]models.TodoDatabaseItem
}

func TestTodo_CreateValidations(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"
	const existingTodoListId = "f2d869a8-e5bc-4fbf-ad71-444444444444"
	const nonExistingTodoListId = "f2d869a8-e5bc-4fbf-ad71-555555555555"

	tests := []createTodoTestCase{
		{
			description:   "Invalid body",
			accessToken:   validAccessToken,
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description: "Access token does not exist",
			accessToken: nonExistingAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id": "%s"}`,
				"fake_description", existingTodoListId),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description: "Description too long",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				strings.Repeat("a", 257), existingTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description: "Description invalid characters",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"/", existingTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description: "Todo list Id invalid",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"description", "not-a-uuid"),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid todo list"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description: "Todo list does not exist",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"description", nonExistingTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"todo list does not exist"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
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
			request.Header.Set("Authorization", tt.accessToken)
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

	tests := []createTodoTestCase{
		{
			description: "Create new todo",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"test todo", existingList),
			responseCode: http.StatusOK,
			responseBody: `{"created_by":"","description":"test todo","status":"todo","updated_at":"2024-06-30T00:00:00+00:00"}`,
			databaseLists: map[string][]models.TodoDatabaseItem{
				existingList: {models.TodoDatabaseItem{
					Description: "test todo",
					Status:      "todo",
					User:        existingAccount,
					UpdatedAt:   fixedTestTime(),
				}},
			},
		},
		{
			description: "Todo list does not exist",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"test todo", nonExistingList),
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"todo list does not exist"}`,
			databaseLists: map[string][]models.TodoDatabaseItem{
				existingList: {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[existingList] = []models.TodoDatabaseItem{}

			todos := Todos{
				Database:     database,
				GenerateUuid: func() string { return "static_uuid" },
				CurrentTime:  func() time.Time { return fixedTestTime() },
			}

			request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(tt.body))
			request.Header.Add("Authorization", tt.accessToken)
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
