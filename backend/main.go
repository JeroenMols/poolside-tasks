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

	users := routes.Users{
		Database:           db.InMemoryDatabase(),
		AddResponseHeaders: net.AddCorsReponseHeaders,
		GenerateUuid:       util.GenerateRandomUuid,
	}

	mux.HandleFunc("POST /users/register", users.Register)
	mux.HandleFunc("POST /users/login", users.Login)
	mux.HandleFunc("GET /users/debug", users.Debug)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}
