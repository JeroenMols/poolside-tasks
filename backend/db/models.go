package db

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	AccountNumber string
	Name          string
}

type AccessToken struct {
	AccountNumber string
	Token         string
}

type TodoList struct {
	Id string
}

type TodoDatabaseItem struct {
	Id          string
	ListId      string
	UpdatedAt   time.Time
	Description string
	Status      string
	User        string
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
