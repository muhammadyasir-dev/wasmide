package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"muhammadyasir-dev/cmd/dbs"
	"muhammadyasir-dev/cmd/models"
)

// Mock database
type MockDB struct {
	CreateFunc func(user *models.User) error
}

func (m *MockDB) Create(user *models.User) error {
	return m.CreateFunc(user)
}

func TestSignup(t *testing.T) {
	// Set up the mock database
	mockDB := &MockDB{
		CreateFunc: func(user *models.User) error {
			// Simulate successful user creation
			return nil
		},
	}

	// Replace the real dbs.Db with the mock
	dbs.Db = mockDB

	tests := []struct {
		name           string
		input          models.User
		expectedStatus int
	}{
		{
			name: "Successful Signup",
			input: models.User{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid JSON",
			input: models.User{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.name == "Invalid JSON" {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			Signup(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.name == "Successful Signup" {
				var user models.User
				if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}

				// Check if the password is hashed
				if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123")); err != nil {
					t.Errorf("password was not hashed correctly")
				}
			}
		})
	}
}
