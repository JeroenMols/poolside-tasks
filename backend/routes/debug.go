package routes

import (
	"backend/db"
	"backend/net"
	"net/http"
)

type Debug struct {
	database *db.Database
}

func CreateDebug(database *db.Database) Debug {
	return Debug{database: database}
}

func (u *Debug) Debug(w http.ResponseWriter, _ *http.Request) {
	net.Success(w, u.database)
}
