package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users/register", func(w http.ResponseWriter, r *http.Request) {
		routeUsersRegister(w)
	})
	mux.HandleFunc("POST /users/login", func(w http.ResponseWriter, r *http.Request) {
		routeUsersLogin(w)
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func routeUsersRegister(w http.ResponseWriter) {
	fmt.Println("Registering new user")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	accountNumber := uuid.New().String()

	w.Write([]byte(fmt.Sprintf("{ \"account_number\": \"%s\" }", accountNumber)))
}

func routeUsersLogin(w http.ResponseWriter) {
	fmt.Println("Logging in new user")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	access_token := uuid.New().String()

	w.Write([]byte(fmt.Sprintf("{ \"access_token\": \"%s\"}", access_token)))
}
