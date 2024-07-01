package routes

import (
	"backend/db"
	"backend/net"
	"fmt"
	"net/http"
)

type TodoLists struct {
	database db.Database
}

func CreateTodoLists(database db.Database) TodoLists {
	return TodoLists{database: database}
}

func (t *TodoLists) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Route todo list Create")

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
	fmt.Printf("Created todo list %s\n", todoList.Id)

	net.Success(w, listCreateResponse{TodoListId: todoList.Id})
}

func (t *TodoLists) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Route todo list Get")
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
		// Ignoring the error, as a real database would handle this using foreign keys
		user, _ := t.database.GetUser(todo.UserId)
		formattedTodos = append(formattedTodos, *toTodoItem(&todo, user))
	}
	fmt.Printf("Get todo list %s\n", listId)

	net.Success(w, listGetResponse{ListId: listId, Todos: formattedTodos})
}
