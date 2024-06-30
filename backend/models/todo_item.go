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
