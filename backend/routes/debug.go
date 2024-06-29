package routes

import (
	"backend/db"
	"backend/net"
	"fmt"
	"net/http"
)

type Debug struct {
	Database           db.Database
	AddResponseHeaders net.AddResponseHeaders
}

func (u *Debug) Debug(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in user")
	u.AddResponseHeaders(w)
	net.Success(w, u.Database)
}
