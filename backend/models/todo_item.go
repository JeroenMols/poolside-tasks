package models

import (
	"backend/db"
	"time"
)

type TodoItem struct {
	Id          string `json:"id"`
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

func ToTodoItem(todo *db.TodoDatabaseItem, user string) TodoItem {
	return TodoItem{
		Id:          todo.Id,
		CreatedBy:   user,
		Description: todo.Description,
		Status:      todo.Status,
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}
}
