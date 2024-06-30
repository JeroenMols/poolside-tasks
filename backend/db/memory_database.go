package db

type Database struct {
	Users        map[string]string
	AccessTokens map[string]string
	TodoLists    map[string][]TodoItem
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]string),
		AccessTokens: make(map[string]string),
		TodoLists:    make(map[string][]TodoItem),
	}
}

// TODO find a better way to put this
type TodoItem struct {
	updatedAt   string
	description string
	status      string
	user        string
}
