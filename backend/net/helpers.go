package net

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Error struct {
	Error string `json:"error,omitempty"`
}

func Success[K any](w http.ResponseWriter, result K) {
	addCorsResponseHeaders(w)
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	_, _ = w.Write(response)
}

func HaltBadRequest(w http.ResponseWriter, error string) {
	addCorsResponseHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	response, _ := json.Marshal(Error{Error: error})
	_, _ = w.Write(response)
}

func ParseBody[K any](r *http.Request) (*K, error) {
	var result K
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&result)
	if err != nil {
		if err.Error() == "unexpected EOF" {
			return nil, errors.New("validation error")
		} else {
			return nil, errors.New("invalid body")
		}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(result)
	if err != nil {
		return nil, errors.New("validation error")
	}

	return &result, nil
}

func addCorsResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
