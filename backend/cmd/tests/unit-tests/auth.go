package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test for loginHandler
func TestLoginHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginHandler)

	// Call the handler with the request and recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	// Check if the redirect URL is correct
	expectedLocation := "http://localhost:5173" // Update this based on your OAuth URL
	if location := rr.Header().Get("Location"); location != expectedLocation {
		t.Errorf("handler returned wrong redirect location: got %v want %v",
			location, expectedLocation)
	}
}

// Test for callbackHandler
func TestCallbackHandler(t *testing.T) {
	// Create a request with a valid state and code
	req, err := http.NewRequest("GET", "/callback?state=validState&code=validCode", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the session and OAuth exchange
	session := &SessionMock{Values: map[string]interface{}{"state": "validState"}}
	store = &StoreMock{session: session}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(callbackHandler)

	// Call the handler with the request and recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	// Check if the redirect URL is correct
	expectedLocation := "http://localhost:5173" // Update this based on your frontend URL
	if location := rr.Header().Get("Location"); location != expectedLocation {
		t.Errorf("handler returned wrong redirect location: got %v want %v",
			location, expectedLocation)
	}
}

// Mock implementations for session and store
type SessionMock struct {
	Values map[string]interface{}
}

func (s *SessionMock) Save(r *http.Request, w http.ResponseWriter) error {
	return nil
}

type StoreMock struct {
	session *SessionMock
}

func (s *StoreMock) Get(r *http.Request, name string) (*SessionMock, error) {
	return s.session, nil
}
