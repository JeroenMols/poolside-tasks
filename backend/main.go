package main

import (
	"backend/db"
	"backend/net"
	"backend/routes"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	database := db.CreateDatabase()

	users := routes.CreateUsers(database)
	todoLists := routes.CreateTodoLists(database)
	todos := routes.CreateTodos(database)

	mux.HandleFunc("POST /users/register", users.Register)
	mux.HandleFunc("POST /users/login", users.Login)

	mux.HandleFunc("POST /todolists", todoLists.Create)
	mux.HandleFunc("GET /todolists/{list_id}", todoLists.Get)

	mux.HandleFunc("POST /todos", todos.Create)
	mux.HandleFunc("PUT /todos/{todo_id}", todos.Update)

	// Debug route
	debug := routes.CreateDebug(&database)
	mux.HandleFunc("GET /debug", debug.Debug)

	authentication := net.AuthenticationMiddleware(mux, database)
	handler := net.CorsMiddleware(authentication, "*")

	fmt.Println("Listening on localhost:8080")
	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		fmt.Println(err.Error())
	}
}
