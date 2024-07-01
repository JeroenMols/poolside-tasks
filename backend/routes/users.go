package routes

import (
	"backend/db"
	"backend/net"
	"fmt"
	"net/http"
	"regexp"
)

type Users struct {
	database db.Database
}

func CreateUsers(database db.Database) Users {
	return Users{database: database}
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new user")

	body, err := net.ParseBody[registerRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	if !regexp.MustCompile(userNameRegex).MatchString(body.Name) {
		net.HaltBadRequest(w, "invalid user name")
		return
	}

	fmt.Printf("User name: %s\n", body.Name)
	user := u.database.CreateUser(body.Name)
	response := registerResponse{
		AccountNumber: user.AccountNumber,
	}

	net.Success(w, response)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")

	body, err := net.ParseBody[loginRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	if !regexp.MustCompile(accountNumberRegex).MatchString(body.AccountNumber) {
		net.HaltBadRequest(w, "invalid account number")
		return
	}

	user, exists := u.database.Users[body.AccountNumber]
	if !exists {
		net.HaltBadRequest(w, "account not found")
		return
	}
	// TODO: return existing token if exists

	fmt.Printf("Account number: %s\n", user.AccountNumber)
	accessToken := u.database.CreateAccessToken(user.AccountNumber)
	response := loginResponse{
		AccessToken: accessToken.Token,
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
