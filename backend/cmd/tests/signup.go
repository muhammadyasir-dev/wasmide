package yourpackage

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"yourmodule/dbs"    // Adjust the import path according to your project structure
	"yourmodule/models" // Adjust the import path according to your project structure
)

func TestSignup(t *testing.T) {
	// Setup a test server
	reqBody := models.User{
		Username: "testuser",
		Password: "testpassword",
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Signup handler
	handler := http.HandlerFunc(Signup)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	var createdUser models.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check if the password is hashed
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("testpassword")); err == nil {
		t.Error("Password should be hashed, but it is not")
	}

	// Check if the username is correct
	if createdUser.Username != reqBody.Username {
		t.Errorf("Expected username %s, got %s", reqBody.Username, createdUser.Username)
	}

	// Clean up: Remove the created user from the database if necessary
	// This part depends on your database setup and should be handled accordingly
}
