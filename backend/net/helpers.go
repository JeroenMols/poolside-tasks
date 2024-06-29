package net

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func Success[K any](w http.ResponseWriter, result K) {
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	_, _ = w.Write(response)
}

func Halt(w http.ResponseWriter, error string) {
	w.WriteHeader(http.StatusBadRequest)
	response, _ := json.Marshal(Error{Error: error})
	_, _ = w.Write(response)
}

func ParseBody[K any](r *http.Request) (K, error) {
	var result K
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&result)
	if err != nil {
		return result, errors.New("invalid body")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(result)
	if err != nil {
		return result, errors.New("validation error")
	}

	return result, nil
}
