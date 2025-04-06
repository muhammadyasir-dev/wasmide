package main_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"muhammadyasir-dev/cmd/apis"
	"muhammadyasir-dev/cmd/dbs"
	"muhammadyasir-dev/cmd/models"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB       *gorm.DB
	testStore    *sessions.CookieStore
	testServer   *httptest.Server
	jwtSecretKey = []byte("test-secret")
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// Initialize test database
	var err error
	testDB, err = gorm.Open(Postgres.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Initialize database schema
	testDB.AutoMigrate(&models.User{}, &models.Fileobject{})
	dbs.Db = testDB

	// Configure OAuth for testing
	apis.OauthConfig = &oauth2.Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"email", "profile"},
	}

	// Configure session store
	testStore = sessions.NewCookieStore([]byte("test-session-key"))
	apis.Store = testStore
	apis.JwtSecret = jwtSecretKey

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/streampty", PsuedoTerminal)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/callback", CallbackHandler)
	mux.HandleFunc("/user", GetUserHandler)
	mux.HandleFunc("/logout", LogoutHandler)

	testServer = httptest.NewServer(mux)
}

func teardown() {
	testServer.Close()
}

func TestAuthFlow(t *testing.T) {
	// Test Signup
	t.Run("Successful user signup", func(t *testing.T) {
		user := map[string]string{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "securepassword",
		}

		body, _ := json.Marshal(user)
		resp, err := http.Post(testServer.URL+"/signup", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdUser models.User
		json.NewDecoder(resp.Body).Decode(&createdUser)
		assert.NotZero(t, createdUser.Id)
		assert.Equal(t, user["name"], createdUser.Name)
	})

	// Test Login and OAuth flow
	t.Run("Complete OAuth flow", func(t *testing.T) {
		// Start login
		resp, err := http.Get(testServer.URL + "/login")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

		// Mock OAuth callback
		req, _ := http.NewRequest("GET", testServer.URL+"/callback?state=test-state&code=test-code", nil)
		req.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		w := httptest.NewRecorder()
		apis.CallbackHandler(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

		// Verify session cookie
		cookies := w.Result().Cookies()
		var authCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "auth_token" {
				authCookie = c
			}
		}
		assert.NotNil(t, authCookie)
	})

	// Test Get User
	t.Run("Get authenticated user", func(t *testing.T) {
		// Create test user and JWT
		user := models.User{
			Name:  "Auth User",
			Email: "auth@example.com",
		}
		testDB.Create(&user)

		token, _ := apis.CreateJWT(apis.UserInfo{Email: user.Email}, user.Id)

		req, _ := http.NewRequest("GET", testServer.URL+"/user", nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var userData models.User
		json.NewDecoder(resp.Body).Decode(&userData)
		assert.Equal(t, user.Email, userData.Email)
	})

	// Test Logout
	t.Run("Successful logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", testServer.URL+"/logout", nil)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookies := resp.Cookies()
		for _, c := range cookies {
			if c.Name == "auth_token" {
				assert.True(t, c.Expires.Before(time.Now()))
			}
		}
	})
}

func TestFileOperations(t *testing.T) {
	// Test Create File
	t.Run("Create new file", func(t *testing.T) {
		fileName := map[string]string{"filename": "test.txt"}
		body, _ := json.Marshal(fileName)
		resp, err := http.Post(testServer.URL+"/create-file", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Test File Upload
	t.Run("Upload file content", func(t *testing.T) {
		content := "test content"
		resp, err := http.Post(testServer.URL+"/files/test.txt", "text/plain", bytes.NewBufferString(content))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Test Get File
	t.Run("Retrieve file content", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/files/test.txt")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "test content")
	})

	// Test List Files
	t.Run("List all files", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/list-files")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var files []string
		json.NewDecoder(resp.Body).Decode(&files)
		assert.Contains(t, files, "test.txt")
	})
}

func TestPseudoTerminal(t *testing.T) {
	t.Run("Execute valid command", func(t *testing.T) {
		resp, err := http.Post(testServer.URL+"/streampty?project=test", "text/plain", bytes.NewBufferString("echo hello"))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "hello")
	})

	t.Run("Handle invalid command", func(t *testing.T) {
		resp, err := http.Post(testServer.URL+"/streampty?project=test", "text/plain", bytes.NewBufferString("invalid-command"))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCORS(t *testing.T) {
	t.Run("CORS preflight request", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", testServer.URL+"/any-endpoint", nil)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	})
}
