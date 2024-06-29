package db

type Database struct {
	Users        map[string]string
	AccessTokens map[string]string
}

func InMemoryDatabase() Database {
	return Database{
		Users:        make(map[string]string),
		AccessTokens: make(map[string]string),
	}
}
