package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreampty(t *testing.T) {
	// Start the server in a separate goroutine
	go main()

	// Wait for the server to start (you may want to implement a better wait mechanism)
	// time.Sleep(time.Second)

	// Prepare the request
	projectName := "testproject"
	command := "echo Hello, World!"
	body, _ := json.Marshal(command)

	req, err := http.NewRequest("POST", "/streampty?project="+projectName, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Streampty)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	assert.Equal(t, "Hello, World!\n", rr.Body.String())
}
