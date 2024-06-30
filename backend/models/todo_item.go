package models

import (
	"time"
)

type TodoDatabaseItem struct {
	UpdatedAt   time.Time
	Description string
	Status      string
	User        string
}

type TodoItem struct {
	CreatedBy   string `json:"created_by"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

func (t *TodoDatabaseItem) ChangeStatus(newStatus string) error {
	if t.Status == "todo" || t.Status == "done" {
		if newStatus == "ongoing" {
			t.Status = newStatus
			return nil
		}
		return errors.New(fmt.Sprintf("invalid status transition from %s to %s", t.Status, newStatus))
	} else if t.Status == "ongoing" {
		if newStatus == "done" || newStatus == "todo" {
			t.Status = newStatus
			return nil
		}
		return errors.New(fmt.Sprintf("invalid status transition from %s to %s", t.Status, newStatus))
	}
	return errors.New(fmt.Sprintf("invalid status transition from %s to %s", t.Status, newStatus))
}

