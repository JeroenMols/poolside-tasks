package net

import "net/http"

type AddResponseHeaders func(w http.ResponseWriter)

func AddCorsReponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
