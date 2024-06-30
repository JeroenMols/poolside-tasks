package routes

import (
	"backend/db"
	"backend/net"
	"backend/util"
	"net/http"
)

type TodoLists struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
}

func (t *TodoLists) Create(w http.ResponseWriter, r *http.Request) {
	net.Success(w, "Create todo list")
}

func (t *TodoLists) Get(w http.ResponseWriter, r *http.Request) {
	listId := r.PathValue("list_id")
	net.Success(w, "Returning todo list "+listId)
}
