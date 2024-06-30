package routes

import (
	"backend/db"
	"backend/net"
	"backend/util"
	"fmt"
	"net/http"
)

type Users struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new user")

	user, err := net.ParseBody[usersRegisterRequest](r)
	if err != nil {
		net.Halt(w, err.Error())
		return
	}

	// TODO: validate user name (forbidden chars, length, etc)

	fmt.Printf("User name: %s\n", user.Name)
	account_number := u.GenerateUuid()
	u.Database.Users[account_number] = user.Name
	response := usersRegisterResponse{
		AccountNumber: account_number,
	}

	net.Success(w, response)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")

	user, err := net.ParseBody[usersLoginRequest](r)
	if err != nil {
		net.Halt(w, err.Error())
		return
	}

	// TODO: check if user exists
	// TODO: return existing token if exists

	fmt.Printf("Account number: %s\n", user.AccountNumber)
	accessToken := u.GenerateUuid()
	u.Database.AccessTokens[user.AccountNumber] = accessToken
	response := usersLoginResponse{
		AccessToken: accessToken,
	}

	net.Success(w, response)
}

type usersRegisterRequest struct {
	Name string `json:"name" validate:"required"`
}

type usersRegisterResponse struct {
	AccountNumber string `json:"account_number"`
}

type usersLoginRequest struct {
	AccountNumber string `json:"account_number" validate:"required"`
}

type usersLoginResponse struct {
	AccessToken string `json:"access_token"`
}
