package db

import (
	"backend/models"
	"backend/util"
	"errors"
	"regexp"
)

type User struct {
	AccountNumber string
	Name          string
}

type TodoList struct {
	Id string
}

type AccessToken struct {
	AccountNumber string
	Token         string
}

type Database struct {
	Users        map[string]User
	AccessTokens map[string]AccessToken
	TodoLists    map[string]TodoList
	TodoItems    map[string]models.TodoDatabaseItem
	currentTime  util.CurrentTime
	generateUuid util.GenerateUuid
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]User),                    // accountNumber -> password
		AccessTokens: make(map[string]AccessToken),             // accessToken -> accountNumber
		TodoLists:    make(map[string]TodoList),                // listId -> listId
		TodoItems:    make(map[string]models.TodoDatabaseItem), // todoId -> todo
		currentTime:  util.GetCurrentTime,
		generateUuid: util.GenerateRandomUuid,
	}
}

func TestDatabase(generateTime util.CurrentTime, generateUuid util.GenerateUuid) Database {
	return Database{
		Users:        make(map[string]User),                    // accountNumber -> password
		AccessTokens: make(map[string]AccessToken),             // accessToken -> accountNumber
		TodoLists:    make(map[string]TodoList),                // listId -> listId
		TodoItems:    make(map[string]models.TodoDatabaseItem), // todoId -> todo
		currentTime:  generateTime,
		generateUuid: generateUuid,
	}
}

const accessTokenRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const listIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const todoIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`

func (d *Database) RegisterUser(name string) string {
	accountNumber := d.generateUuid()
	d.Users[accountNumber] = User{AccountNumber: accountNumber, Name: name}
	return accountNumber
}

func (d *Database) CreateAccessToken(accountNumber string) *AccessToken {
	accessToken := AccessToken{AccountNumber: accountNumber, Token: d.generateUuid()}
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
	todoList := TodoList{Id: d.generateUuid()}
	d.TodoLists[todoList.Id] = todoList
	return &todoList
}

func (d *Database) CreateTodo(listId string, description string, user string) models.TodoDatabaseItem {
	itemId := d.generateUuid()
	item := models.TodoDatabaseItem{
		Id:          itemId,
		ListId:      listId,
		Description: description,
		Status:      "todo",
		User:        user,
		UpdatedAt:   d.currentTime(),
	}

	d.TodoItems[item.Id] = item
	return item
}

func (d *Database) UpdateTodo(todo *models.TodoDatabaseItem) error {
	_, exists := d.TodoItems[todo.Id]
	if !exists {
		return errors.New("todo does not exist")
	}
	todo.UpdatedAt = d.currentTime()
	d.TodoItems[todo.Id] = *todo
	return nil
}

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

func (d *Database) GetTodos(listId string) (*[]models.TodoDatabaseItem, error) {
	if !regexp.MustCompile(listIdRegex).MatchString(listId) {
		return nil, errors.New("invalid todo list")
	}

	todoList, exists := d.TodoLists[listId]
	if !exists {
		return nil, errors.New("todo list does not exist")
	}

	var items []models.TodoDatabaseItem
	for _, item := range d.TodoItems {
		if item.ListId == todoList.Id {
			items = append(items, item)
		}
	}

	return &items, nil
}
