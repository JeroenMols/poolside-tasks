package routes

import (
	"backend/net"
	"backend/util"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type Users struct {
	AddResponseHeaders net.AddResponseHeaders
	GenerateUuid       util.GenerateUuid
}

type Error struct {
	Error string `json:"error"`
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new user")
	u.AddResponseHeaders(w)

	if r.Body == nil {
		return
	}

	var user usersRegisterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.Marshal(Error{Error: "Validation error"})
		if err != nil {
			_, _ = w.Write([]byte("{\"error\":\"Invalid body\"}"))
		}
		fmt.Println(err.Error())
		return
	} else {
		fmt.Printf("User name: %s\n", user.Name)
	}

	response := usersRegisterResponse{
		AccountNumber: u.GenerateUuid(),
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
