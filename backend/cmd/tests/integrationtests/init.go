package apis

import (
	"os"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestInit(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("SESSION_KEY", "test-session-key")
	os.Setenv("JWT_SECRET", "test-jwt-key")
	os.Setenv("CLIENT_ID", "test-client-id")
	os.Setenv("CLIENT_SECRET", "test-client-secret")
	os.Setenv("REDIRECT_URL", "http://localhost:8080/auth/callback")

	// Call the init function
	init()

	// Validate session store
	assert.NotNil(t, store, "Session store should be initialized")
	assert.Equal(t, "test-session-key", string(store.Options.Path), "Session path should match")

	// Validate JWT secret
	assert.Equal(t, "test-jwt-key", string(jwtSecret), "JWT secret should match")

	// Validate OAuth configuration
	assert.NotNil(t, oauthConfig, "OAuth configuration should be initialized")
	assert.Equal(t, "test-client-id", oauthConfig.ClientID, "Client ID should match")
	assert.Equal(t, "test-client-secret", oauthConfig.ClientSecret, "Client Secret should match")
	assert.Equal(t, "http://localhost:8080/auth/callback", oauthConfig.RedirectURL, "Redirect URL should match")
	assert.Contains(t, oauthConfig.Scopes, "https://www.googleapis.com/auth/userinfo.email", "Scopes should include userinfo.email")
	assert.Contains(t, oauthConfig.Scopes, "https://www.googleapis.com/auth/userinfo.profile", "Scopes should include userinfo.profile")
}

func TestGetEnvWithDefault(t *testing.T) {
	// Test with an existing environment variable
	os.Setenv("TEST_VAR", "test-value")
	result := getEnvWithDefault("TEST_VAR", "default-value")
	assert.Equal(t, "test-value", result, "Should return the value of the environment variable")

	// Test with a non-existing environment variable
	result = getEnvWithDefault("NON_EXISTING_VAR", "default-value")
	assert.Equal(t, "default-value", result, "Should return the default value")
}
