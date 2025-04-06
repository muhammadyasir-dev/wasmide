package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var testDB *gorm.DB

// Setup function to initialize the in-memory database
func setup() {
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	testDB.AutoMigrate(&User {})
	db = testDB // Use the test database
}

// TestCreateUser  tests the CreateUser Handler
func TestCreateUser (t *testing.T) {
	setup()
	defer teardown()

	user := User{Name: "John Doe", Email: "john@example.com"}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser Handler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdUser  User
	json.NewDecoder(rr.Body).Decode(&createdUser )
	assert.Equal(t, user.Name, createdUser .Name)
	assert.Equal(t, user.Email, createdUser .Email)
	assert.NotZero(t, createdUser .ID) // Ensure ID is set
}

// TestGetUser  tests the GetUser Handler
func TestGetUser (t *testing.T) {
	setup()
	defer teardown()

	// Create a user first
	user := User{Name: "Jane Doe", Email: "jane@example.com"}
	testDB.Create(&user)

	req, err := http.NewRequest("GET", "/users/"+string(user.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser Handler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedUser  User
	json.NewDecoder(rr.Body).Decode(&retrievedUser )
	assert.Equal(t, user.Name, retrievedUser .Name)
	assert.Equal(t, user.Email, retrievedUser .Email)
}

// TestUpdateUser  tests the UpdateUser Handler
func TestUpdateUser (t *testing.T) {
	setup()
	defer teardown()

	// Create a user first
	user := User{Name: "Alice", Email: "alice@example.com"}
	testDB.Create(&user)

	updatedUser  := User{Name: "Alice Updated", Email: "alice_updated@example.com"}
	body, _ := json.Marshal(updatedUser )

	req, err := http.NewRequest("PUT", "/users/"+string(user.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateUser Handler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedUser  User
	json.NewDecoder(rr.Body).Decode(&retrievedUser )
	assert.Equal(t, updatedUser .Name, retrievedUser .Name)
	assert.Equal(t, updatedUser .Email, retrievedUser .Email)
}

// TestDeleteUser  tests the DeleteUser Handler
func TestDeleteUser (t *testing.T) {
	setup()
	defer teardown()

	// Create a user first
	user := User{Name: "Bob", Email: "bob@example.com"}
	testDB.Create(&user)

	req, err := http.NewRequest("DELETE", "/users/"+string(user.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteUser Handler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify the user is deleted
	var deletedUser  User
	result := testDB.First(&deletedUser , user.ID)
	assert.Error(t, result.Error) // Expect an error since the user should not exist
}

// Teardown function to clean up the database
func teardown() {
	testDB.Exec("DELETE FROM users")
}
