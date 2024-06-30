package db

import (
	"backend/models"
	"errors"
	"regexp"
)

type Database struct {
	Users        map[string]string
	AccessTokens map[string]string
	TodoLists    map[string][]models.TodoDatabaseItem
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]string),
		AccessTokens: make(map[string]string),
		TodoLists:    make(map[string][]models.TodoDatabaseItem),
	}
}

const accessTokenRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const listIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const todoIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`

func (d *Database) Authorize(accessToken string) (*string, error) {
	if !regexp.MustCompile(accessTokenRegex).MatchString(accessToken) {
		return nil, errors.New("invalid access token")
	}
	accountNumber := d.AccessTokens[accessToken]
	if accountNumber == "" {
		return nil, errors.New("account not found")
	}
	return &accountNumber, nil
}

func (d *Database) GetTodos(listId string) (*[]models.TodoDatabaseItem, error) {
	if !regexp.MustCompile(listIdRegex).MatchString(listId) {
		return nil, errors.New("invalid todo list")
	}

	todoList := d.TodoLists[listId]
	if todoList == nil {
		return nil, errors.New("todo list does not exist")
	}
	return &todoList, nil
}

func (d *Database) GetTodo(todoId string) (*models.TodoDatabaseItem, error) {
	if !regexp.MustCompile(todoIdRegex).MatchString(todoId) {
		return nil, errors.New("invalid todo")
	}

	todo, err := d.findTodo(todoId)
	if err != nil {
		return nil, errors.New("todo does not exist")
	}
	return todo, nil
}

// Note: this would not be required when using a database, as todo items would
// live in their own table and we could query based on the todoId.
func (d *Database) findTodo(todoId string) (*models.TodoDatabaseItem, error) {
	for _, todos := range d.TodoLists {
		var todo models.TodoDatabaseItem
		for _, todo = range todos {
			if todo.Id == todoId {
				return &todo, nil
			}
		}
	}
	return nil, errors.New("todo not found")
}
