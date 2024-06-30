package routes

import (
	"backend/db"
	"backend/models"
	"backend/net"
	"backend/util"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Todos struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
	CurrentTime  util.CurrentTime
}

func (t *Todos) Create(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[todoCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}
	accountNumber, err := t.Database.Authorize(body.AccessToken)
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	if !regexp.MustCompile(todoDescriptionRegex).MatchString(body.Description) {
		net.HaltBadRequest(w, "invalid description")
		return
	}

	if !regexp.MustCompile(listIdRegex).MatchString(body.ListId) {
		net.HaltBadRequest(w, "invalid todo list")
		return
	}

	todoList := t.Database.TodoLists[body.ListId]
	if todoList == nil {
		net.HaltBadRequest(w, "todo list does not exist")
		return
	}

	item := models.TodoItem{
		Description: body.Description,
		Status:      "todo",
		User:        *accountNumber,
		UpdatedAt:   t.CurrentTime(),
	}
	fmt.Printf("Creating new todo %s\n", item.Description)
	t.Database.TodoLists[body.ListId] = append(todoList, item)

	net.Success(w, todoCreateResponse{
		CreatedBy:   t.Database.Users[*accountNumber],
		Description: item.Description,
		Status:      item.Status,
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	})
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
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}
