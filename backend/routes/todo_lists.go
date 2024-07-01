package routes

import (
	"backend/db"
	"backend/net"
	"fmt"
	"net/http"
	"time"
)

type TodoLists struct {
	database db.Database
}

func CreateTodoLists(database db.Database) TodoLists {
	return TodoLists{database: database}
}

func (t *TodoLists) Create(w http.ResponseWriter, r *http.Request) {
	_, err := net.ParseBody[listCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}
	if _, err := t.database.GetAccessToken(r.Header.Get("Authorization")); err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	todoList := t.database.CreateTodoList()
	fmt.Printf("Creating new todo list %s\n", todoList.Id)

	net.Success(w, listCreateResponse{TodoListId: todoList.Id})
}

func (t *TodoLists) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := t.database.GetAccessToken(r.Header.Get("Authorization")); err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	listId := r.PathValue("list_id")

	todos, err := t.database.GetTodos(listId)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	formattedTodos := []todoItem{}
	for _, todo := range *todos {
		formattedTodos = append(formattedTodos, todoItem{
			Id:          todo.Id,
			CreatedBy:   t.database.Users[todo.UserId].Name,
			Description: todo.Description,
			Status:      todo.Status,
			UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
		})
	}

	net.Success(w, listGetResponse{
		ListId: listId,
		Todos:  formattedTodos,
	})
}
