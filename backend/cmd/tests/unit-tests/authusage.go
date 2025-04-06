package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestMain runs first to load environment variables
func TestMain(m *testing.M) {
	// Setup: Set test environment variables
	os.Setenv("SESSION_KEY", "test-session-key")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("CLIENT_ID", "test-client-id")
	os.Setenv("CLIENT_SECRET", "test-client-secret")

	// Run init() manually after setting env vars
	init()

	// Run tests
	code := m.Run()

	// Teardown
	os.Unsetenv("SESSION_KEY")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("CLIENT_ID")
	os.Unsetenv("CLIENT_SECRET")

	os.Exit(code)
}

func TestGetEnvWithDefault(t *testing.T) {
	t.Run("returns environment value when set", func(t *testing.T) {
		os.Setenv("TEST_KEY", "expected-value")
		defer os.Unsetenv("TEST_KEY")

		result := getEnvWithDefault("TEST_KEY", "default")
		assert.Equal(t, "expected-value", result)
	})

	t.Run("returns default value when not set", func(t *testing.T) {
		result := getEnvWithDefault("NON_EXISTENT_KEY", "default-value")
		assert.Equal(t, "default-value", result)
	})
}

func TestRoutesRegistration(t *testing.T) {
	router := mux.NewRouter()
	setupRoutes(router)

	tests := []struct {
		path   string
		method string
	}{
		{"/auth/login", "GET"},
		{"/auth/callback", "GET"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			route := router.Get(tt.path)
			assert.NotNil(t, route, "Route not registered: %s", tt.path)

			methods, err := route.GetMethods()
			assert.NoError(t, err)
			assert.Contains(t, methods, tt.method)
		})
	}
}

func TestCORSConfiguration(t *testing.T) {
	router := mux.NewRouter()
	setupRoutes(router)
	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(router)

	t.Run("preflight request returns CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/auth/login", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		req.Header.Set("Access-Control-Request-Method", "GET")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "http://localhost:5173", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET,POST,OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	})
}

func TestSessionStoreConfiguration(t *testing.T) {
	assert.NotNil(t, store, "Session store not initialized")
	assert.Equal(t, 604800, store.Options.MaxAge, "Session max age should be 7 days")
	assert.True(t, store.Options.HttpOnly, "Session should be HTTP only")
	assert.Equal(t, http.SameSiteLaxMode, store.Options.SameSite)
}

func TestOAuthConfigInitialization(t *testing.T) {
	assert.Equal(t, "test-client-id", oauthConfig.ClientID)
	assert.Equal(t, "test-client-secret", oauthConfig.ClientSecret)
	assert.Equal(t, "http://localhost:8080/auth/callback", oauthConfig.RedirectURL)
	assert.Equal(t, google.Endpoint, oauthConfig.Endpoint)
}

// Helper to replicate route setup
func setupRoutes(r *mux.Router) {
	r.HandleFunc("/auth/login", loginHandler).Methods("GET")
	r.HandleFunc("/auth/callback", callbackHandler).Methods("GET")
}
