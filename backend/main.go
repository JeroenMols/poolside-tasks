package main

import (
	"backend/net"
	"backend/routes"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	users := routes.Users{
		AddResponseHeaders: net.AddCorsReponseHeaders,
	}

	mux.HandleFunc("POST /users/register", users.Register)
	mux.HandleFunc("POST /users/login", users.Login)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}
