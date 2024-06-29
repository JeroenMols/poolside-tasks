package main

import (
	"backend/db"
	"backend/net"
	"backend/routes"
	"backend/util"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	database := db.InMemoryDatabase()

	users := routes.Users{
		Database:           database,
		AddResponseHeaders: net.AddCorsReponseHeaders,
		GenerateUuid:       util.GenerateRandomUuid,
	}
	todoLists := routes.TodoLists{
		Database:           database,
		AddResponseHeaders: net.AddCorsReponseHeaders,
		GenerateUuid:       util.GenerateRandomUuid,
	}

	mux.HandleFunc("POST /users/register", users.Register)
	mux.HandleFunc("POST /users/login", users.Login)
	mux.HandleFunc("POST /todolists", todoLists.Create)
	mux.HandleFunc("GET /todolists/{list_id}", todoLists.Get)

	// Debug route
	debug := routes.Debug{
		Database:           database,
		AddResponseHeaders: net.AddCorsReponseHeaders,
	}
	mux.HandleFunc("GET /debug", debug.Debug)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}
