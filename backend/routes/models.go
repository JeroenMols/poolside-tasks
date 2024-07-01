package routes

import (
	"backend/db"
	"time"
)

type registerRequest struct {
	Name string `json:"name" validate:"required"`
}

type registerResponse struct {
	UserId string `json:"user_id"`
}

type loginRequest struct {
	UserId string `json:"user_id" validate:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

type listCreateRequest struct {
}

type listCreateResponse struct {
	TodoListId string `json:"todo_list_id"`
}

type listGetResponse struct {
	ListId string     `json:"todo_list_id"`
	Todos  []todoItem `json:"todos"`
}

type todoCreateRequest struct {
	ListId      string `json:"todo_list_id" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type todoUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}

type todoItem struct {
	Id          string `json:"id"`
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

func toTodoItem(todo *db.TodoItem, user *db.User) *todoItem {
	return &todoItem{
		Id:          todo.Id,
		CreatedBy:   user.Name,
		Description: todo.Description,
		Status:      todo.Status,
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}
}
