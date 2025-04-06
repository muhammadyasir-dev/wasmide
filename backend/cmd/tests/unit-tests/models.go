package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"muhammadyasir-dev/cmd/models"
)

func TestPsuedoTerminal(t *testing.T) {
	req := httptest.NewRequest("GET", "/pseudo-terminal", nil)
	w := httptest.NewRecorder()

	PsuedoTerminal(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestSignup(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Signup(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	// Optionally, you can check the response body if needed
}

func TestLoginHandler(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	// Optionally, you can check the response body if needed
}

func TestCallbackHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/callback", nil)
	w := httptest.NewRecorder()

	CallbackHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestGetUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	GetUserHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestLogoutHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()

	LogoutHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}
