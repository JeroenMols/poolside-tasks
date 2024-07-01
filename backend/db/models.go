package db

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id   string
	Name string
}

type AccessToken struct {
	UserId string
	Token  string
}

type TodoList struct {
	Id string
}

type TodoItem struct {
	Id          string
	ListId      string
	UserId      string
	UpdatedAt   time.Time
	Description string
	Status      string
}

func (t *TodoItem) ChangeStatus(newStatus string) error {
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
