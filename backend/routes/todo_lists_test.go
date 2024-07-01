package routes

import (
	"backend/db"
	"backend/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type createListTestCase struct {
	description   string
	accessToken   string
	body          string
	responseCode  int
	responseBody  string
	databaseLists map[string]db.TodoList
}

func TestTodoLists_Create(t *testing.T) {
	tests := []createListTestCase{
		{
			description:   "Invalid body",
			accessToken:   fakeToken,
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"body not valid"}`,
			databaseLists: make(map[string]db.TodoList),
		},
		{
			description:   "Access token does not exist",
			accessToken:   fakeWrongToken,
			body:          `{}`,
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"user not found"}`,
			databaseLists: make(map[string]db.TodoList),
		},
		{
			description:   "Create new todo list",
			accessToken:   fakeToken,
			body:          `{}`,
			responseCode:  http.StatusOK,
			responseBody:  `{"todo_list_id":"static_uuid"}`,
			databaseLists: map[string]db.TodoList{"static_uuid": {Id: "static_uuid"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2021, 1, 1) },
				func(string) string { return "static_uuid" },
			)
			database.AccessTokens[fakeToken] = db.AccessToken{UserId: fakeUserId, Token: fakeToken}

			todoList := CreateTodoLists(database)

			request := httptest.NewRequest(http.MethodPost, "/todolists", strings.NewReader(tt.body))
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todoList.Create(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseLists, database.TodoLists)
		})
	}
}

type getListTestCase struct {
	description  string
	accessToken  string
	todoListId   string
	responseCode int
	responseBody string
}

func TestTodoLists_Get(t *testing.T) {
	const fakeNoElementsTodoListId = fakeTodoListId2

	tests := []getListTestCase{
		{
			description:  "Invalid access token",
			accessToken:  "not-an-access-token",
			todoListId:   fakeNoElementsTodoListId,
			responseCode: http.StatusUnauthorized,
			responseBody: `{"error":"invalid access token"}`,
		},
		{
			description:  "Invalid todo list id parameter",
			accessToken:  fakeToken,
			todoListId:   `invalid-list-id`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid todo list"}`,
		},
		{
			description:  "todo list not found",
			accessToken:  fakeToken,
			todoListId:   fakeWrongTodoListId,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"todo list not found"}`,
		},
		{
			description:  "Get empty todo list",
			accessToken:  fakeToken,
			todoListId:   fakeNoElementsTodoListId,
			responseCode: http.StatusOK,
			responseBody: `{"todo_list_id":"lst_cccccccccccccccccccccc","todos":[]}`,
		},
		{
			description:  "Get todo list",
			accessToken:  fakeToken,
			todoListId:   fakeTodoListId,
			responseCode: http.StatusOK,
			responseBody: `{"todo_list_id":"lst_aaaaaaaaaaaaaaaaaaaaaa","todos":[{"id":"id1","created_by":"test user","description":"first todo","status":"todo","updated_at":"2024-01-01T00:00:00+00:00"},{"id":"id2","created_by":"test user","description":"second todo","status":"ongoing","updated_at":"2023-01-01T00:00:00+00:00"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.TestDatabase(
				func() time.Time { return util.FakeTime(2021, 1, 1) },
				func(string) string { return "static_uuid" },
			)
			database.Users[fakeUserId] = db.User{Id: fakeUserId, Name: "test user"}
			database.AccessTokens[fakeToken] = db.AccessToken{UserId: fakeUserId, Token: fakeToken}
			database.TodoLists[fakeNoElementsTodoListId] = db.TodoList{Id: fakeNoElementsTodoListId}
			database.TodoLists[fakeTodoListId] = db.TodoList{Id: fakeTodoListId}
			database.TodoItems = map[string]db.TodoItem{
				"id1": {Id: "id1", ListId: fakeTodoListId, Description: "first todo", Status: "todo", UserId: fakeUserId, UpdatedAt: util.FakeTime(2024, 1, 1)},
				"id2": {Id: "id2", ListId: fakeTodoListId, Description: "second todo", Status: "ongoing", UserId: fakeUserId, UpdatedAt: util.FakeTime(2023, 1, 1)},
			}
			database.TodoItemOrder = []string{"id1", "id2"}

			todoList := CreateTodoLists(database)

			request := httptest.NewRequest(http.MethodGet, "/todolists", nil)
			request.SetPathValue("list_id", tt.todoListId)
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todoList.Get(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
		})
	}
}
