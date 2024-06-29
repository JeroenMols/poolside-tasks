package routes

import (
	"backend/db"
	"backend/net"
	"backend/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Users struct {
	Database           db.Database
	AddResponseHeaders net.AddResponseHeaders
	GenerateUuid       util.GenerateUuid
}

type Error struct {
	Error string `json:"error"`
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new user")
	u.AddResponseHeaders(w)

	user, err := parseBody[usersRegisterRequest](r)
	if err != nil {
		halt(w, err.Error())
		return
	}

	fmt.Printf("User name: %s\n", user.Name)
	account_number := u.GenerateUuid()
	u.Database.Users[account_number] = user.Name
	response := usersRegisterResponse{
		AccountNumber: account_number,
	}

	success(w, response)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")
	u.AddResponseHeaders(w)

	user, err := parseBody[usersLoginRequest](r)
	if err != nil {
		halt(w, err.Error())
		return
	}

	fmt.Printf("Account number: %s\n", user.AccountNumber)
	accessToken := u.GenerateUuid()
	u.Database.AccessTokens[user.AccountNumber] = accessToken
	response := usersLoginResponse{
		AccessToken: accessToken,
	}

	success(w, response)
}

func (u *Users) Debug(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")
	u.AddResponseHeaders(w)

	success(w, u.Database)
}

func parseBody[K any](r *http.Request) (K, error) {
	var result K
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&result)
	if err != nil {
		return result, errors.New("invalid body")
	}
	return result, nil
}

func success[K any](w http.ResponseWriter, result K) {
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	_, _ = w.Write(response)
}

func halt(w http.ResponseWriter, error string) {
	w.WriteHeader(http.StatusBadRequest)
	response, _ := json.Marshal(Error{Error: error})
	_, _ = w.Write(response)
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
