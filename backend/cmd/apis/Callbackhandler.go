/*

package apis

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"io/ioutil"
	"muhammadyasir-dev/cmd/dbs"
	"muhammadyasir-dev/cmd/models"
	"net/http"
	"time"
)

// UserInfo represents the user information from Google
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

var (
	oauthConfig *oauth2.Config
	store       *sessions.CookieStore
	jwtSecret   []byte
)

func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

	// Log the exact URL we're redirecting to
	url := oauthConfig.AuthCodeURL(state)
	fmt.Println("Redirecting to OAuth URL:", url)

	// Redirect to Google's OAuth page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
	fmt.Printf(":0x99]   Received cookies: %+v\n", r.Cookies())
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

	// Print user info to console for debugging
	fmt.Printf("User  logged in: %+v\n", userInfo)

	dbUser, err := findOrCreateUser(userInfo)
	if err != nil {
		http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Create JWT token
	jwtToken, err := CreateJWT(userInfo, User.ID)
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
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	// Redirect to frontend

	fmt.Printf(":0x145]   Received cookies: %+v\n", r.Cookies())

	frontendURL := "http://localhost:5173"
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func findaOrCreateUser(UserInfo UserInfo) (*User, error) {
	var user models.User

	// Try to find user by email first
	result := dbs.Db.Where("email = ?", userInfo.Email).First(&user)
	if result.Error == nil {
		// User found, return
		return &user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		// Some other error occurred
		return nil, result.Error
	}

	// Try to find user by name as fallback (if you need this)
	result = dbs.Db.Where("name = ?", userInfo.Name).First(&user)
	if result.Error == nil {
		// User found by name, return
		return &user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		// Some other error occurred
		return nil, result.Error
	}

	// User not found, create new user
	newUser := User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Picture:  userInfo.Picture,
		GoogleID: userInfo.ID,
		// Password field is left empty for OAuth users
	}

	// Save user to database
	result = dbs.Db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func CreateJWT(userInfo UserInfo) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"id":      userInfo.ID,
		"email":   userInfo.Email,
		"name":    userInfo.Name,
		"picture": userInfo.Picture,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires in 7 days
		"iat":     time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString(jwtSecret)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get JWT from cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse and validate JWT
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	// Extract user data from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusInternalServerError)
		return
	}

	// Prepare user data
	userData := map[string]interface{}{
		"id":      claims["id"],
		"email":   claims["email"],
		"name":    claims["name"],
		"picture": claims["picture"],
	}

	// Send user data as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)

	// Also log user data to console
	fmt.Printf("User  data from JWT: %+v\n", userData)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear auth cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}


*/

package apis

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
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
	Password string `gorm:"column:password;default:''" json:"password,omitempty"`
	Picture  string `gorm:"column:picture" json:"picture,omitempty"`
	GoogleID string `gorm:"column:google_id" json:"google_id,omitempty"`
	// You can add more fields as needed
}

// UserInfo represents the user information from Google
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
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
	// Auto migrate the User model
	db.AutoMigrate(&User{})
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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

	// Log the exact URL we're redirecting to
	url := oauthConfig.AuthCodeURL(state)
	fmt.Println("Redirecting to OAuth URL:", url)

	// Redirect to Google's OAuth page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

	fmt.Printf("Received cookies: %+v\n", r.Cookies())

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

	// Print user info to console for debugging
	fmt.Printf("User logged in: %+v\n", userInfo)

	// Check if user exists in database and create if not
	dbUser, err := findOrCreateUser(userInfo)
	if err != nil {
		http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create JWT token
	jwtToken, err := CreateJWT(userInfo, dbUser.ID)
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
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	// Redirect to frontend
	fmt.Printf("Received cookies: %+v\n", r.Cookies())
	frontendURL := "http://localhost:5173"
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

// findOrCreateUser checks if a user exists in the database and creates one if not
func findOrCreateUser(userInfo UserInfo) (*User, error) {
	var user User

	// Try to find user by email first
	result := db.Where("email = ?", userInfo.Email).First(&user)
	if result.Error == nil {
		// User found, return
		return &user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		// Some other error occurred
		return nil, result.Error
	}

	// Try to find user by name as fallback (if you need this)
	result = db.Where("name = ?", userInfo.Name).First(&user)
	if result.Error == nil {
		// User found by name, return
		return &user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		// Some other error occurred
		return nil, result.Error
	}

	// User not found, create new user
	newUser := User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Picture:  userInfo.Picture,
		GoogleID: userInfo.ID,
		// Password field is left empty for OAuth users
	}

	// Save user to database
	result = db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func CreateJWT(userInfo UserInfo, userID uint) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"id":        userID,
		"google_id": userInfo.ID,
		"email":     userInfo.Email,
		"name":      userInfo.Name,
		"picture":   userInfo.Picture,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires in 7 days
		"iat":       time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString(jwtSecret)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get JWT from cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse and validate JWT
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	// Extract user data from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusInternalServerError)
		return
	}

	// Get user from database
	var user User
	result := db.First(&user, claims["id"])
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Send user data as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	// Also log user data to console
	fmt.Printf("User data from database: %+v\n", user)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear auth cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
