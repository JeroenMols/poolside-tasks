package db

import (
	"backend/util"
	"errors"
	"regexp"
)

type Database struct {
	Users        map[string]User
	AccessTokens map[string]AccessToken
	TodoLists    map[string]TodoList
	TodoItems    map[string]TodoDatabaseItem
	currentTime  util.CurrentTime
	generateUuid util.GenerateUuid
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]User),
		AccessTokens: make(map[string]AccessToken),
		TodoLists:    make(map[string]TodoList),
		TodoItems:    make(map[string]TodoDatabaseItem),
		currentTime:  util.GetCurrentTime,
		generateUuid: util.GenerateRandomUuid,
	}
}

func TestDatabase(generateTime util.CurrentTime, generateUuid util.GenerateUuid) Database {
	return Database{
		Users:        make(map[string]User),
		AccessTokens: make(map[string]AccessToken),
		TodoLists:    make(map[string]TodoList),
		TodoItems:    make(map[string]TodoDatabaseItem),
		currentTime:  generateTime,
		generateUuid: generateUuid,
	}
}

const accessTokenRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const listIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const todoIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`

func (d *Database) CreateUser(name string) *User {
	user := User{
		AccountNumber: d.generateUuid(),
		Name:          name,
	}
	d.Users[user.AccountNumber] = user
	return &user
}

func (d *Database) CreateAccessToken(accountNumber string) *AccessToken {
	accessToken := AccessToken{
		AccountNumber: accountNumber,
		Token:         d.generateUuid(),
	}
	d.AccessTokens[accessToken.Token] = accessToken
	return &accessToken
}

func (d *Database) GetAccessToken(token string) (*AccessToken, error) {
	if !regexp.MustCompile(accessTokenRegex).MatchString(token) {
		return nil, errors.New("invalid access token")
	}
	accessToken, exists := d.AccessTokens[token]
	if !exists {
		return nil, errors.New("account not found")
	}
	return &accessToken, nil
}

func (d *Database) CreateTodoList() *TodoList {
	todoList := TodoList{
		Id: d.generateUuid(),
	}
	d.TodoLists[todoList.Id] = todoList
	return &todoList
}

func (d *Database) CreateTodo(listId string, description string, user string) *TodoDatabaseItem {
	item := TodoDatabaseItem{
		Id:          d.generateUuid(),
		ListId:      listId,
		Description: description,
		Status:      "todo",
		User:        user,
		UpdatedAt:   d.currentTime(),
	}

	d.TodoItems[item.Id] = item
	return &item
}

// TODO: should also return todo here
func (d *Database) UpdateTodo(todo *TodoDatabaseItem) error {
	_, exists := d.TodoItems[todo.Id]
	if !exists {
		return errors.New("todo does not exist")
	}
	todo.UpdatedAt = d.currentTime()
	d.TodoItems[todo.Id] = *todo
	return nil
}

func (d *Database) GetTodo(todoId string) (*TodoDatabaseItem, error) {
	if !regexp.MustCompile(todoIdRegex).MatchString(todoId) {
		return nil, errors.New("invalid todo")
	}

	item, exists := d.TodoItems[todoId]
	if !exists {
		return nil, errors.New("todo does not exist")
	}
	return &item, nil
}

func (d *Database) GetTodos(listId string) (*[]TodoDatabaseItem, error) {
	if !regexp.MustCompile(listIdRegex).MatchString(listId) {
		return nil, errors.New("invalid todo list")
	}

	todoList, exists := d.TodoLists[listId]
	if !exists {
		return nil, errors.New("todo list does not exist")
	}

	var items []TodoDatabaseItem
	for _, item := range d.TodoItems {
		if item.ListId == todoList.Id {
			items = append(items, item)
		}
	}

	return &items, nil
}
