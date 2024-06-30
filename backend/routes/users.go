package routes

import (
	"backend/db"
	"backend/net"
	"backend/util"
	"fmt"
	"net/http"
	"regexp"
)

type Users struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new user")

	user, err := net.ParseBody[registerRequest](r)
	if err != nil {
		net.Halt(w, err.Error())
		return
	}

	if !regexp.MustCompile(userNameRegex).MatchString(user.Name) {
		net.Halt(w, "invalid user name")
		return
	}

	fmt.Printf("User name: %s\n", user.Name)
	accountNumber := u.GenerateUuid()
	u.Database.Users[accountNumber] = user.Name
	response := registerResponse{
		AccountNumber: accountNumber,
	}

	net.Success(w, response)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")

	user, err := net.ParseBody[loginRequest](r)
	if err != nil {
		net.Halt(w, err.Error())
		return
	}

	if !regexp.MustCompile(accountNumberRegex).MatchString(user.AccountNumber) {
		net.Halt(w, "invalid account number")
		return
	}

	if u.Database.Users[user.AccountNumber] == "" {
		net.Halt(w, "account not found")
		return
	}
	// TODO: return existing token if exists

	fmt.Printf("Account number: %s\n", user.AccountNumber)
	accessToken := u.GenerateUuid()
	u.Database.AccessTokens[user.AccountNumber] = accessToken
	response := loginResponse{
		AccessToken: accessToken,
	}

	net.Success(w, response)
}

type registerRequest struct {
	Name string `json:"name" validate:"required"`
}

type registerResponse struct {
	AccountNumber string `json:"account_number"`
}

type loginRequest struct {
	AccountNumber string `json:"account_number" validate:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}
