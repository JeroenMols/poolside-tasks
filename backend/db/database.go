package db

import (
	"backend/models"
	"errors"
	"regexp"
)

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

const accessTokenRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`

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
