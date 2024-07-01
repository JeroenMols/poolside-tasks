package routes

import (
	"backend/db"
	"backend/net"
	"net/http"
)

type Debug struct {
	Database db.Database
}

func (u *Debug) Debug(w http.ResponseWriter, _ *http.Request) {
	net.Success(w, u.Database)
}
