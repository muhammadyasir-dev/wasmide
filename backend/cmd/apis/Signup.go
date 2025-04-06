package apis

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"muhammadyasir-dev/cmd/dbs"
	"muhammadyasir-dev/cmd/models"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var signupuser models.User

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&signupuser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupuser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	signupuser.Password = string(hashedPassword) // Store the hashed password

	// Log the database operations
	// Create the user in the database
	if err := dbs.Db.Create(&signupuser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signupuser)
}
