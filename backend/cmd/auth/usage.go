// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
	store       *sessions.CookieStore
	jwtSecret   []byte
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Initialize session store with a secure key
	sessionKey := []byte(getEnvWithDefault("SESSION_KEY", "super-secret-session-key"))
	store = sessions.NewCookieStore(sessionKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	// JWT Secret
	jwtSecret = []byte(getEnvWithDefault("JWT_SECRET", "super-secret-jwt-key"))

	// OAuth Configuration
	redirectURL := getEnvWithDefault("REDIRECT_URL", "http://localhost:8080/auth/callback")
	oauthConfig = &oauth2.Config{
		ClientID:     getEnvWithDefault("CLIENT_ID", ""),
		ClientSecret: getEnvWithDefault("CLIENT_SECRET", ""),
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	// Validate essential configuration
	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		log.Println("Warning: CLIENT_ID or CLIENT_SECRET not set. OAuth will not work correctly.")
	}
}

func main() {
	r := mux.NewRouter()

	// Auth endpoints
	r.HandleFunc("/auth/login", loginHandler).Methods("GET")
	r.HandleFunc("/auth/callback", callbackHandler).Methods("GET")

	// Other routes...

	// CORS configuration
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Vite's default port
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	// Wrap the router with CORS middleware
	handler := corsHandler(r)

	// Print the routes for debugging
	fmt.Println("Available Routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		fmt.Printf("Path: %s, Methods: %v\n", path, methods)
		return nil
	})

	fmt.Println("Server started at http://localhost:8080")

	// Make double sure we're using the router
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
