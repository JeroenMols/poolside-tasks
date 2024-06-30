package routes

import (
	"backend/db"
	"backend/models"
	"backend/net"
	"backend/util"
	"fmt"
	"net/http"
	"regexp"
)

type Todos struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
}

func (t *Todos) Create(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[todoCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	if !regexp.MustCompile(accessTokenRegex).MatchString(body.AccessToken) {
		net.HaltUnauthorized(w, "invalid access token")
		return
	}

	accountNumber := t.Database.AccessTokens[body.AccessToken]
	if accountNumber == "" {
		net.HaltUnauthorized(w, "account not found")
		return
	}

	if !regexp.MustCompile(todoDescriptionRegex).MatchString(body.Description) {
		net.HaltBadRequest(w, "invalid description")
		return
	}

	listId := t.GenerateUuid()
	fmt.Printf("Creating new todo list %s\n", listId)
	t.Database.TodoLists[listId] = []db.TodoItem{}

	net.Success(w, listCreateResponse{TodoListId: listId})
}

func (t *Todos) Update(w http.ResponseWriter, r *http.Request) {
	net.Success(w, "Updated TODO")
}

type todoCreateRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
	Description string `json:"description" validate:"required"`
	ListId      string `json:"todo_list_id" validate:"required"`
}

type todoCreateResponse struct {
}
