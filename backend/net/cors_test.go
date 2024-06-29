package net

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCorsReponseHeaders(t *testing.T) {
	w := httptest.NewRecorder()

	AddCorsReponseHeaders(w)

	assert.Equal(t, http.Header{"Access-Control-Allow-Origin": []string{"*"}}, w.Result().Header)
}
