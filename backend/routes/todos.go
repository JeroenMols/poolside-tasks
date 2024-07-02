package routes

import (
	"backend/db"
	"backend/net"
	"fmt"
	"net/http"
	"regexp"
)

type Todos struct {
	database db.Database
}

func CreateTodos(database db.Database) Todos {
	return Todos{database: database}
}

func (t *Todos) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Route todo Create")
	body, err := net.ParseBody[todoCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	accessToken, err := t.database.GetAccessToken(r.Header.Get("Authorization"))
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	if !regexp.MustCompile(todoDescriptionRegex).MatchString(body.Description) {
		net.HaltBadRequest(w, "description not valid")
		return
	}

	_, err = t.database.GetTodos(body.ListId)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	item := t.database.CreateTodo(body.ListId, body.Description, accessToken.UserId)

	// No need to handle error, we already know the user exists
	user, _ := t.database.GetUser(item.UserId)

	fmt.Printf("Created todo %s\n", item.Id)

	net.Success(w, toTodoItem(item, user))
}

func (t *Todos) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Route todo Update")
	body, err := net.ParseBody[todoUpdateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	todoId := r.PathValue("todo_id")

	_, err = t.database.GetAccessToken(r.Header.Get("Authorization"))
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	item, err := t.database.GetTodo(todoId)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	err = item.ChangeStatus(body.Status)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	// No need to handle error, we already know the both exists
	updatedItem, _ := t.database.UpdateTodo(item)
	user, _ := t.database.GetUser(updatedItem.UserId)
	fmt.Printf("Updated todo %s\n", item.Id)

	net.Success(w, toTodoItem(updatedItem, user))
}
