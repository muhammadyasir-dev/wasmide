package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

// User represents the user model in the database
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Email    string `gorm:"column:email;uniqueIndex" json:"email"`
	Picture  string `gorm:"column:picture" json:"picture,omitempty"`
	GoogleID string `gorm:"column:google_id" json:"google_id,omitempty"`
}

// UserInfo represents the user information from Google
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

var (
	oauthConfig *oauth2.Config
	store       *sessions.CookieStore
	jwtSecret   []byte
	db          *gorm.DB
)

// InitDB initializes the database connection
func InitDB(database *gorm.DB) {
	db = database
	db.AutoMigrate(&User {})
}

func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	state, err := GenerateStateToken()
	if err != nil {
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
		return
	}

	// Store state in session
	session, _ := store.Get(r, "oauth-state")
	session.Values["state"] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to Google's OAuth page
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state
	session, _ := store.Get(r, "oauth-state")
	expectedState, ok := session.Values["state"].(string)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	if state != expectedState {
		http.Error(w, fmt.Sprintf("Invalid state parameter. Expected: %s, Got: %s", expectedState, state), http.StatusBadRequest)
		return
	}

	// Exchange code for token
	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info
	client := oauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var userInfo UserInfo
	if err = json.Unmarshal(data, &userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user exists in database and create if not
	dbUser , err := findOrCreateUser (userInfo)
	if err != nil {
		http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create JWT token
	jwtToken, err := CreateJWT(userInfo, dbUser .ID)
	if err != nil {
		http.Error(w, "Failed to create JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set JWT in cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, cookie)

	// Redirect to frontend
	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

// findOrCreateUser  checks if a user exists in the database and creates one if not
func findOrCreateUser (userInfo UserInfo) (*User , error) {
	var user User

	// Try to find user by email first
	result := db.Where("email = ?", userInfo.Email).First(&user)
	if result.Error == nil {
		return &user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	// User not found, create new user
	newUser  := User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Picture:  userInfo.Picture,
		GoogleID: userInfo.ID,
	}

	// Save user to database
	result = db.Create(&newUser )
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser , nil
}

func CreateJWT(userInfo UserInfo, userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":        userID,
		"google_id": userInfo.ID,
		"email":     userInfo.Email,
		"name":      userInfo.Name,
		"picture":   userInfo.Picture,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetUser Handler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusInternalServerError)
		return
	}

	var user User
	result := db.First(&user, claims["id"])
	if result.Error != nil {
		http.Error(w, "User  not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie :=










