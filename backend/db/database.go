package db

import (
	"backend/models"
	"errors"
	"regexp"
	"time"
)

type Database struct {
	Users        map[string]string
	AccessTokens map[string]string
	TodoLists    map[string]string
	TodoItems    map[string]models.TodoDatabaseItem
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]string),                  // accountNumber -> password
		AccessTokens: make(map[string]string),                  // accessToken -> accountNumber
		TodoLists:    make(map[string]string),                  // listId -> listId
		TodoItems:    make(map[string]models.TodoDatabaseItem), // todoId -> todo
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

// TODO list id should become internal
func (d *Database) CreateTodoList(uuid string) string {
	listId := uuid
	d.TodoLists[listId] = listId
	return listId
}

// TODO todo id and time should become internal
func (d *Database) CreateTodo(uuid string, listId string, description string, user string, updatedAt time.Time) models.TodoDatabaseItem {
	item := models.TodoDatabaseItem{
		Id:          uuid,
		ListId:      listId,
		Description: description,
		Status:      "todo",
		User:        user,
		UpdatedAt:   updatedAt,
	}

	d.TodoItems[item.Id] = item
	return item
}

// TODO might need explicit tests
func (d *Database) GetTodos(listId string) (*[]models.TodoDatabaseItem, error) {
	if !regexp.MustCompile(listIdRegex).MatchString(listId) {
		return nil, errors.New("invalid todo list")
	}

	todoList, exists := d.TodoLists[listId]
	if !exists {
		return nil, errors.New("todo list does not exist")
	}

	items := make([]models.TodoDatabaseItem, 0, len(todoList))
	for _, item := range d.TodoItems {
		if item.ListId == listId {
			items = append(items, item)
		}
	}

	return &items, nil
}

// TODO might need explicit tests
func (d *Database) GetTodo(todoId string) (*models.TodoDatabaseItem, error) {
	if !regexp.MustCompile(todoIdRegex).MatchString(todoId) {
		return nil, errors.New("invalid todo")
	}

	item, exists := d.TodoItems[todoId]
	if !exists {
		return nil, errors.New("todo does not exist")
	}
	return &item, nil
}

// TODO might need explicit tests
func (d *Database) UpdateTodo(todo *models.TodoDatabaseItem) error {
	_, exists := d.TodoItems[todo.Id]
	if !exists {
		return errors.New("todo does not exist")
	}
	d.TodoItems[todo.Id] = *todo
	return nil
}
