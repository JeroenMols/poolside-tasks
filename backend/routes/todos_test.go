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

type createTodoTestCase struct {
	description   string
	accessToken   string
	body          string
	responseCode  int
	responseBody  string
	databaseTodos map[string]db.TodoItem
}

func TestTodo_CreateValidations(t *testing.T) {
	tests := []createTodoTestCase{
		{
			description:   "Invalid body",
			accessToken:   fakeToken,
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
		{
			description: "Access token does not exist",
			accessToken: fakeWrongToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id": "%s"}`,
				"fake_description", fakeTodoListId),
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
		{
			description: "Description too long",
			accessToken: fakeToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				strings.Repeat("a", 257), fakeTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
		{
			description: "Description invalid characters",
			accessToken: fakeToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"/", fakeTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid description"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
		{
			description: "Todo list Id invalid",
			accessToken: fakeToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"description", "not-a-uuid"),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid todo list"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
		{
			description: "Todo list does not exist",
			accessToken: fakeToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"description", fakeWrongTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"todo list does not exist"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func(string) string { return "static_uuid" },
			)
			database.AccessTokens[fakeToken] = db.AccessToken{UserId: fakeUserId, Token: fakeToken}
			database.TodoLists[fakeTodoListId] = db.TodoList{Id: fakeTodoListId}

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
	tests := []createTodoTestCase{
		{
			description: "Create new todo",
			accessToken: fakeToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"test todo", fakeTodoListId),
			responseCode: http.StatusOK,
			responseBody: `{"id":"static_uuid","created_by":"","description":"test todo","status":"todo","updated_at":"2024-06-30T00:00:00+00:00"}`,
			databaseTodos: map[string]db.TodoItem{
				"static_uuid": {
					Id:          "static_uuid",
					ListId:      fakeTodoListId,
					Description: "test todo",
					Status:      "todo",
					UserId:      fakeUserId,
					UpdatedAt:   util.FakeTime(2024, 6, 30),
				},
			},
		},
		{
			description: "Todo list does not exist",
			accessToken: fakeToken,
			body: fmt.Sprintf(`{"description":"%s", "todo_list_id":"%s"}`,
				"test todo", fakeWrongTodoListId),
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"todo list does not exist"}`,
			databaseTodos: make(map[string]db.TodoItem),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func(string) string { return "static_uuid" },
			)
			database.AccessTokens[fakeToken] = db.AccessToken{UserId: fakeUserId, Token: fakeToken}
			database.TodoLists[fakeTodoListId] = db.TodoList{Id: fakeTodoListId}

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
	databaseLists map[string][]db.TodoItem
}

func TestTodo_UpdateValidations(t *testing.T) {
	tests := []updateTodoTestCase{
		{
			description:  "Invalid body",
			accessToken:  fakeToken,
			todoId:       fakeTodoId,
			body:         `{"invalid":"body"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid body"}`,
		},
		{
			description:  "Access token does not exist",
			accessToken:  fakeWrongToken,
			todoId:       fakeTodoId,
			body:         `{"status":"progress"}`,
			responseCode: http.StatusUnauthorized,
			responseBody: `{"error":"account not found"}`,
		},
		{
			description:  "Todo Id invalid",
			accessToken:  fakeToken,
			todoId:       "not-a-uuid",
			body:         `{"status":"progress"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid todo"}`,
		},
		{
			description:  "Todo does not exist",
			accessToken:  fakeToken,
			todoId:       fakeWrongTodoId,
			body:         `{"status":"progress"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"todo does not exist"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func(string) string { return "static_uuid" },
			)
			database.AccessTokens[fakeToken] = db.AccessToken{UserId: fakeUserId, Token: fakeToken}
			database.TodoLists[fakeTodoListId] = db.TodoList{Id: fakeTodoListId}
			database.TodoItems = map[string]db.TodoItem{
				fakeTodoId: {Id: fakeTodoId, ListId: fakeTodoListId, Description: "first todo", Status: "todo", UserId: fakeUserId, UpdatedAt: util.FakeTime(2024, 1, 1)},
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
	tests := []updateTodoTestCase{
		{
			description:  "Update todo valid transition",
			accessToken:  fakeToken,
			todoId:       fakeTodoId,
			body:         `{"status":"ongoing"}`,
			responseCode: http.StatusOK,
			responseBody: `{"id":"tdo_aaaaaaaaaaaaaaaaaaaaaa","created_by":"","description":"first todo","status":"ongoing","updated_at":"2024-06-30T00:00:00+00:00"}`,
			databaseLists: map[string][]db.TodoItem{
				fakeTodoListId: {db.TodoItem{
					Id:          "static_uuid",
					Description: "first todo",
					Status:      "ongoing",
					UserId:      fakeUserId,
					UpdatedAt:   util.FakeTime(2024, 6, 30),
				}},
			},
		},
		{
			description:  "Update todo invalid transition",
			accessToken:  fakeToken,
			todoId:       fakeTodoId,
			body:         `{"status":"done"}`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid status transition from todo to done"}`,
			databaseLists: map[string][]db.TodoItem{
				fakeTodoListId: {db.TodoItem{
					Id:          "static_uuid",
					Description: "first todo",
					Status:      "ongoing",
					UserId:      fakeUserId,
					UpdatedAt:   util.FakeTime(2000, 1, 1),
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2024, 6, 30) },
				func(string) string { return "static_uuid" },
			)
			database.AccessTokens[fakeToken] = db.AccessToken{UserId: fakeUserId, Token: fakeToken}
			database.TodoLists[fakeTodoListId] = db.TodoList{Id: fakeTodoListId}
			database.TodoItems = map[string]db.TodoItem{
				fakeTodoId: {Id: fakeTodoId, ListId: fakeTodoListId, Description: "first todo", Status: "todo", UserId: fakeUserId, UpdatedAt: util.FakeTime(2000, 1, 1)},
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
