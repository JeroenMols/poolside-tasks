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

	fmt.Printf("UserId name: %s\n", body.Name)
	user := u.database.CreateUser(body.Name)
	response := registerResponse{
		UserId: user.Id,
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

	if !regexp.MustCompile(accountNumberRegex).MatchString(body.UserId) {
		net.HaltBadRequest(w, "invalid user id")
		return
	}

	user, exists := u.database.Users[body.UserId]
	if !exists {
		net.HaltBadRequest(w, "user not found")
		return
	}
	// TODO: return existing token if exists

	fmt.Printf("Account number: %s\n", user.Id)
	accessToken := u.database.CreateAccessToken(user.Id)
	response := loginResponse{
		AccessToken: accessToken.Token,
	}

	net.Success(w, response)
}
