package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStreampty(t *testing.T) {
	// Create a test server
	handler := http.HandlerFunc(Streampty)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Define the project name and command to execute
	projectName := "test-project"
	command := "echo Hello, World!"

	// Prepare the request body
	reqBody := bytes.NewBufferString(command)

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?project=%s", server.URL, projectName), reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Send the request to the server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got: %s, body: %s", resp.Status, body)
	}

	// Check the output
	expectedOutput := "Hello, World!\n"
	if string(body) != expectedOutput {
		t.Fatalf("Expected output %q, got %q", expectedOutput, string(body))
	}
}

func TestMain(m *testing.M) {
	// Setup code can go here if needed

	// Run the tests
	code := m.Run()

	// Cleanup code can go here if needed

	// Exit with the code from the tests
	os.Exit(code)
}
