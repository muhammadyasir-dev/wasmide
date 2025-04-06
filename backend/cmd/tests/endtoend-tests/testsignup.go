package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents the user model in the database
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Database instance
var db *gorm.DB

// Initialize the database
func initDB() {
	var err error
	db, err = gorm.Open(Postgres.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&User{})
}

// Signup Handler handles user registration
func Signup(w http.ResponseWriter, r *http.Request) {
	var signupUser User

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&signupUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	signupUser.Password = string(hashedPassword) // Store the hashed password

	// Create the user in the database
	if err := db.Create(&signupUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signupUser)
}

// Main function to start the server
func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/signup", Signup).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
