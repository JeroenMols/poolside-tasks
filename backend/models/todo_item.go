package models

import (
	"errors"
	"fmt"
	"time"
)

type TodoDatabaseItem struct {
	Id          string
	ListId      string
	UpdatedAt   time.Time
	Description string
	Status      string
	User        string
}

type TodoItem struct {
	Id          string `json:"id"`
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

func (t *TodoDatabaseItem) ChangeStatus(newStatus string, time time.Time) error {
	if t.Status == "todo" || t.Status == "done" {
		if newStatus == "ongoing" {
			t.Status = newStatus
			t.UpdatedAt = time
			return nil
		}
		return errors.New(fmt.Sprintf("invalid status transition from %s to %s", t.Status, newStatus))
	} else if t.Status == "ongoing" {
		if newStatus == "done" || newStatus == "todo" {
			t.Status = newStatus
			t.UpdatedAt = time
			return nil
		}
		return errors.New(fmt.Sprintf("invalid status transition from %s to %s", t.Status, newStatus))
	}
	return errors.New(fmt.Sprintf("invalid status transition from %s to %s", t.Status, newStatus))
}

func (t *TodoDatabaseItem) ToTodoItem(user string) TodoItem {
	return TodoItem{
		Id:          t.Id,
		CreatedBy:   user,
		Description: t.Description,
		Status:      t.Status,
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}
