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
	TodoItems    map[string]TodoItem
	currentTime  util.CurrentTime
	generateUuid util.GenerateUuid
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]User),
		AccessTokens: make(map[string]AccessToken),
		TodoLists:    make(map[string]TodoList),
		TodoItems:    make(map[string]TodoItem),
		currentTime:  util.GetCurrentTime,
		generateUuid: util.GenerateRandomUuid,
	}
}

func TestDatabase(generateTime util.CurrentTime, generateUuid util.GenerateUuid) Database {
	return Database{
		Users:        make(map[string]User),
		AccessTokens: make(map[string]AccessToken),
		TodoLists:    make(map[string]TodoList),
		TodoItems:    make(map[string]TodoItem),
		currentTime:  generateTime,
		generateUuid: generateUuid,
	}
}

const accessTokenRegex = `^tkn_[23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy=]{22}$`
const listIdRegex = `^lst_[23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy=]{22}$`
const todoIdRegex = `^tdo_[23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy=]{22}$`

func (d *Database) CreateUser(name string) *User {
	user := User{
		Id:   d.generateUuid("usr"),
		Name: name,
	}
	d.Users[user.Id] = user
	return &user
}

func (d *Database) GetUser(userId string) (*User, error) {
	user, exists := d.Users[userId]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (d *Database) CreateAccessToken(accountNumber string) *AccessToken {
	accessToken := AccessToken{
		UserId: accountNumber,
		Token:  d.generateUuid("tkn"),
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
		return nil, errors.New("user not found")
	}
	return &accessToken, nil
}

func (d *Database) CreateTodoList() *TodoList {
	todoList := TodoList{
		Id: d.generateUuid("lst"),
	}
	d.TodoLists[todoList.Id] = todoList
	return &todoList
}

func (d *Database) CreateTodo(listId string, description string, user string) *TodoItem {
	item := TodoItem{
		Id:          d.generateUuid("tdo"),
		ListId:      listId,
		Description: description,
		Status:      "todo",
		UserId:      user,
		UpdatedAt:   d.currentTime(),
	}

	d.TodoItems[item.Id] = item
	return &item
}

func (d *Database) UpdateTodo(todo *TodoItem) (*TodoItem, error) {
	_, exists := d.TodoItems[todo.Id]
	if !exists {
		return nil, errors.New("todo not found")
	}
	todo.UpdatedAt = d.currentTime()
	d.TodoItems[todo.Id] = *todo
	return todo, nil
}

func (d *Database) GetTodo(todoId string) (*TodoItem, error) {
	if !regexp.MustCompile(todoIdRegex).MatchString(todoId) {
		return nil, errors.New("invalid todo")
	}

	item, exists := d.TodoItems[todoId]
	if !exists {
		return nil, errors.New("todo not found")
	}
	return &item, nil
}

func (d *Database) GetTodos(listId string) (*[]TodoItem, error) {
	if !regexp.MustCompile(listIdRegex).MatchString(listId) {
		return nil, errors.New("invalid todo list")
	}

	todoList, exists := d.TodoLists[listId]
	if !exists {
		return nil, errors.New("todo list not found")
	}

	var items []TodoItem
	for _, item := range d.TodoItems {
		if item.ListId == todoList.Id {
			items = append(items, item)
		}
	}

	return &items, nil
}
