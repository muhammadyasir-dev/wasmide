package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"muhammadyasir-dev/cmd/handler"
)

func TestRouter(t *testing.T) {
	router := Router()

	t.Run("Test Signup", func(t *testing.T) {
		user := map[string]string{"username": "testuser", "password": "testpass"}
		body, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Signup handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Test Login", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/login", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Login handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Test Get User", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/user", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GetUser  handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Test Logout", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/logout", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Logout handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Test Auth Callback", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/auth/callback", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Callback handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Test Stream", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/stream", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("PsuedoTerminal handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
