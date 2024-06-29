package routes

import (
	"backend/net"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type Users struct {
	AddResponseHeaders net.AddResponseHeaders
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new user")
	u.AddResponseHeaders(w)

	var user usersRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("User name: %s\n", user.Name)
	}

	accountNumber := uuid.New().String()
	response := usersRegisterResponse{
		AccountNumber: accountNumber,
	}

	responseString, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Write(responseString)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")
	u.AddResponseHeaders(w)

	var user usersLoginRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Account number: %s\n", user.AccountNumber)
	}

	accessToken := uuid.New().String()
	response := usersLoginResponse{
		AccessToken: accessToken,
	}

	responseString, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Write(responseString)
}

type usersRegisterRequest struct {
	Name string `json:"name"`
}

type usersRegisterResponse struct {
	AccountNumber string `json:"account_number"`
}

type usersLoginRequest struct {
	AccountNumber string `json:"account_number"`
}

type usersLoginResponse struct {
	AccessToken string `json:"access_token"`
}
