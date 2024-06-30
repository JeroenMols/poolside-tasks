package db

import "backend/models"

type Database struct {
	Users        map[string]string
	AccessTokens map[string]string
	TodoLists    map[string][]models.TodoItem
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]string),
		AccessTokens: make(map[string]string),
		TodoLists:    make(map[string][]models.TodoItem),
	}
}
