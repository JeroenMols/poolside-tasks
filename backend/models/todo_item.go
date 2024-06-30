package models

import (
	"time"
)

type TodoItem struct {
	UpdatedAt   time.Time
	Description string
	Status      string
	User        string
}
