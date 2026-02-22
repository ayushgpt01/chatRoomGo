package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	testData := User{ID: 1, Name: "Jane Doe"}
	expectedStatus := http.StatusCreated
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	err := Encode(w, r, expectedStatus, testData)

	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content-type application/json, got %s", contentType)
	}

	var actualData User
	if err := json.NewDecoder(w.Body).Decode(&actualData); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if actualData != testData {
		t.Errorf("expected body %v, got %v", testData, actualData)
	}
}

type MockValidator struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (m MockValidator) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if m.Name == "" {
		problems["name"] = "name is required"
	}
	if m.Age < 18 {
		problems["age"] = "must be an adult"
	}
	return problems
}

func TestDecodeValid(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantProblems int
		wantErr      bool
	}{
		{
			name:         "Valid Input",
			body:         `{"name": "Alice", "age": 30}`,
			wantProblems: 0,
			wantErr:      false,
		},
		{
			name:         "Validation Failure",
			body:         `{"name": "", "age": 10}`,
			wantProblems: 2,
			wantErr:      true,
		},
		{
			name:         "Malformed JSON",
			body:         `{"name": "Alice", "age": "not-a-number"}`,
			wantProblems: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request with string body
			r := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))

			// Call the generic function with our mock type
			val, problems, err := DecodeValid[MockValidator](r)

			// Assertions
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeValid() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(problems) != tt.wantProblems {
				t.Errorf("expected %d problems, got %d", tt.wantProblems, len(problems))
			}

			if !tt.wantErr && val.Name == "" {
				t.Error("expected name to be populated in valid case")
			}
		})
	}
}
