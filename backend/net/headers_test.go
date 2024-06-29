package net

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCorsResponseHeaders(t *testing.T) {
	w := httptest.NewRecorder()

	AddCorsReponseHeaders(w)

	assert.Equal(t, http.Header{
		"Content-Type":                []string{"application/json"},
		"Access-Control-Allow-Origin": []string{"*"},
	}, w.Result().Header)
}
