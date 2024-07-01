package routes

import (
	"backend/db"
	"backend/models"
	"backend/util"
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
	databaseTodos map[string]models.TodoDatabaseItem
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
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
		{
			description: "Access token does not exist",
			accessToken: nonExistingAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id": "%s"}`,
				"fake_description", existingTodoListId),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
		{
			description: "Description too long",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				strings.Repeat("a", 257), existingTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
		{
			description: "Description invalid characters",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"/", existingTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
		{
			description: "Todo list Id invalid",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"description", "not-a-uuid"),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid todo list"}`,
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
		{
			description: "Todo list does not exist",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"description", nonExistingTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"todo list does not exist"}`,
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func() string { return "static_uuid" },
			)
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[existingTodoListId] = db.TodoList{Id: existingTodoListId}

			todos := CreateTodos(database)

			request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(tt.body))
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todos.Create(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseTodos, database.TodoItems)
		})
	}
}

func TestTodo_CreateLogic(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const existingList = "f2d869a8-e5bc-4fbf-ad71-444444444444"
	const nonExistingList = "f2d869a8-e5bc-4fbf-ad71-555555555555"

	tests := []createTodoTestCase{
		{
			description: "Create new todo",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"test todo", existingList),
			responseCode: http.StatusOK,
			responseBody: `{"id":"static_uuid","created_by":"","description":"test todo","status":"todo","updated_at":"2024-06-30T00:00:00+00:00"}`,
			databaseTodos: map[string]models.TodoDatabaseItem{
				"static_uuid": {
					Id:          "static_uuid",
					ListId:      existingList,
					Description: "test todo",
					Status:      "todo",
					User:        existingAccount,
					UpdatedAt:   util.FakeTime(2024, 6, 30),
				},
			},
		},
		{
			description: "Todo list does not exist",
			accessToken: validAccessToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"test todo", nonExistingList),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"todo list does not exist"}`,
			databaseTodos: make(map[string]models.TodoDatabaseItem),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func() string { return "static_uuid" },
			)
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[existingList] = db.TodoList{Id: existingList}

			todos := CreateTodos(database)

			request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(tt.body))
			request.Header.Add("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todos.Create(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseTodos, database.TodoItems)
		})
	}
}

type updateTodoTestCase struct {
	description   string
	accessToken   string
	todoId        string
	body          string
	responseCode  int
	responseBody  string
	databaseLists map[string][]models.TodoDatabaseItem
}

func TestTodo_UpdateValidations(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"
	const existingList = "f2d869a8-e5bc-4fbf-ad71-444444444444"
	const nonExistingList = "f2d869a8-e5bc-4fbf-ad71-555555555555"
	const existingTodoId = "f2d869a8-e5bc-4fbf-ad71-6666666666666"
	const nonExistingTodoId = "f2d869a8-e5bc-4fbf-ad71-777777777777"

	tests := []updateTodoTestCase{
		{
			description:  "Invalid body",
			accessToken:  validAccessToken,
			todoId:       existingTodoId,
			body:         `{"invalid":"body"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid body"}`,
		},
		{
			description:  "Access token does not exist",
			accessToken:  nonExistingAccessToken,
			todoId:       existingTodoId,
			body:         `{"status":"progress"}`,
			responseCode: http.StatusUnauthorized,
			responseBody: `{"error":"account not found"}`,
		},
		{
			description:  "Todo Id invalid",
			accessToken:  validAccessToken,
			todoId:       "not-a-uuid",
			body:         `{"status":"progress"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid todo"}`,
		},
		{
			description:  "Todo does not exist",
			accessToken:  validAccessToken,
			todoId:       nonExistingTodoId,
			body:         `{"status":"progress"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"todo does not exist"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func() string { return "static_uuid" },
			)
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[existingList] = db.TodoList{Id: existingList}
			database.TodoItems = map[string]models.TodoDatabaseItem{
				existingTodoId: {Id: existingTodoId, ListId: existingList, Description: "first todo", Status: "todo", User: existingAccount, UpdatedAt: util.FakeTime(2024, 1, 1)},
			}

			todos := CreateTodos(database)

			request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(tt.body))
			request.SetPathValue("todo_id", tt.todoId)
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todos.Update(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
		})
	}
}

func TestTodos_UpdateLogic(t *testing.T) {
	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"
	const existingList = "f2d869a8-e5bc-4fbf-ad71-444444444444"
	const nonExistingList = "f2d869a8-e5bc-4fbf-ad71-555555555555"
	const existingTodoId = "f2d869a8-e5bc-4fbf-ad71-666666666666"
	const nonExistingTodoId = "f2d869a8-e5bc-4fbf-ad71-777777777777"

	tests := []updateTodoTestCase{
		{
			description:  "Update todo valid transition",
			accessToken:  validAccessToken,
			todoId:       existingTodoId,
			body:         `{"status":"ongoing"}`,
			responseCode: http.StatusOK,
			responseBody: `{"id":"f2d869a8-e5bc-4fbf-ad71-666666666666","created_by":"","description":"first todo","status":"ongoing","updated_at":"2024-06-30T00:00:00+00:00"}`,
			databaseLists: map[string][]models.TodoDatabaseItem{
				existingList: {models.TodoDatabaseItem{
					Id:          "static_uuid",
					Description: "first todo",
					Status:      "ongoing",
					User:        existingAccount,
					UpdatedAt:   util.FakeTime(2024, 6, 30),
				}},
			},
		},
		{
			description:  "Update todo invalid transition",
			accessToken:  validAccessToken,
			todoId:       existingTodoId,
			body:         `{"status":"done"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid status transition from todo to done"}`,
			databaseLists: map[string][]models.TodoDatabaseItem{
				existingList: {models.TodoDatabaseItem{
					Id:          "static_uuid",
					Description: "first todo",
					Status:      "ongoing",
					User:        existingAccount,
					UpdatedAt:   util.FakeTime(2000, 1, 1),
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func() string { return "static_uuid" },
			)
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[existingList] = db.TodoList{Id: existingList}
			database.TodoItems = map[string]models.TodoDatabaseItem{
				existingTodoId: {Id: existingTodoId, ListId: existingList, Description: "first todo", Status: "todo", User: existingAccount, UpdatedAt: util.FakeTime(2000, 1, 1)},
			}

			todos := CreateTodos(database)

			request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(tt.body))
			request.SetPathValue("todo_id", tt.todoId)
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todos.Update(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
		})
	}
}
