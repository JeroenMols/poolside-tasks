package net

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testData struct {
	Name string `json:"name,omitempty" validate:"required"`
	Age  int    `json:"age,omitempty" validate:"required"`
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		input  testData
		output string
	}{
		{testData{Name: "test"}, "{\"name\":\"test\"}"},
		{testData{}, "{}"},
		{testData{Name: ""}, "{}"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Scenario %d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			Success(w, tt.input)
			assert.Equal(t, 200, w.Result().StatusCode)
			assert.Equal(t, tt.output, w.Body.String())
		})
	}
}

func TestHaltBadRequest(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"body validation error", "{\"error\":\"body validation error\"}"},
		{"", "{}"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Scenario %d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			HaltBadRequest(w, tt.input)
			assert.Equal(t, 400, w.Result().StatusCode)
			assert.Equal(t, tt.output, w.Body.String())
		})
	}
}

func TestHaltUnauthorized(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"body validation error", "{\"error\":\"body validation error\"}"},
		{"", "{}"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Scenario %d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			HaltUnauthorized(w, tt.input)
			assert.Equal(t, 401, w.Result().StatusCode)
			assert.Equal(t, tt.output, w.Body.String())
		})
	}
}

type parseBodyTestCase struct {
	description string
	body        string
	output      *testData
	err         string
}

func TestParseBody(t *testing.T) {
	tests := []parseBodyTestCase{
		{
			description: "Missing body",
			body:        "",
			output:      nil,
			err:         "body not valid",
		},
		{
			description: "Invalid body",
			body:        `body not valid`,
			output:      nil,
			err:         "body not valid",
		},
		{
			description: "Empty body",
			body:        "{}",
			output:      nil,
			err:         "validation error",
		},
		{
			description: "Missing attribute",
			body:        `{"name":"test"}`,
			output:      nil,
			err:         "validation error",
		},
		{
			description: "Additional attribute",
			body:        `{"name":"test", "age": 30, "invalid":"attribute}`,
			output:      nil,
			err:         "validation error",
		},
		{
			description: "Valid body",
			body:        `{"name":"test", "age": 30}`,
			output:      &testData{Name: "test", Age: 30},
			err:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			result, err := ParseBody[testData](r)

			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Equal(t, err, nil)
			}

			assert.Equal(t, tt.output, result)
		})
	}
}
