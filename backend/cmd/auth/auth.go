// auth.go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	state, err := generateStateToken()
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

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Log to help debug - which path was called?
	fmt.Printf("Callback received on path: %s\n", r.URL.Path)

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

	// Print user info to console for debugging
	fmt.Printf("User  logged in: %+v\n", userInfo)

	// Create JWT token
	jwtToken, err := createJWT(userInfo)
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
	frontendURL := "http://localhost:5173"
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
