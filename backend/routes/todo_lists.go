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

type TodoLists struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
}

func (t *TodoLists) Create(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[listCreateRequest](r)
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

	listId := t.GenerateUuid()
	fmt.Printf("Creating new todo list %s\n", listId)
	t.Database.TodoLists[listId] = []models.TodoItem{}

	net.Success(w, listCreateResponse{TodoListId: listId})
}

func (t *TodoLists) Get(w http.ResponseWriter, r *http.Request) {
	listId := r.PathValue("list_id")
	net.Success(w, "Returning todo list "+listId)
}

type listCreateRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type listCreateResponse struct {
	TodoListId string `json:"todo_list_id"`
}
