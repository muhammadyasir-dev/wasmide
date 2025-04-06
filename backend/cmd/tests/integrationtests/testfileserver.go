package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testFileDir = "./test_files"

func setup() {
	// Create a test directory for files
	os.RemoveAll(testFileDir) // Clean up before test
	os.Mkdir(testFileDir, os.ModePerm)
	fileDir = testFileDir // Use the test directory
}

func teardown() {
	// Clean up after tests
	os.RemoveAll(testFileDir)
}

func TestCreateFileHandler(t *testing.T) {
	setup()
	defer teardown()

	// Create a new file
	fileName := "testfile.txt"
	body, _ := json.Marshal(fileName)
	req, err := http.NewRequest("POST", "/create-file", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&Server{}).createFileHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("CreateFileHandler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check if the file was created
	if _, err := os.Stat(filepath.Join(fileDir, fileName)); os.IsNotExist(err) {
		t.Errorf("File %s was not created", fileName)
	}
}

func TestListFilesHandler(t *testing.T) {
	setup()
	defer teardown()

	// Create a new file to list
	fileName := "testfile.txt"
	body, _ := json.Marshal(fileName)
	req, err := http.NewRequest("POST", "/create-file", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&Server{}).createFileHandler)
	handler.ServeHTTP(rr, req)

	// Now list files
	req, err = http.NewRequest("GET", "/list-files", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc((&Server{}).listFilesHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ListFilesHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var files []string
	if err := json.NewDecoder(rr.Body).Decode(&files); err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 || files[0] != fileName {
		t.Errorf("ListFilesHandler returned unexpected files: got %v want %v", files, []string{fileName})
	}
}

func TestFileHandler(t *testing.T) {
	setup()
	defer teardown()

	// Create a new file
	fileName := "testfile.txt"
	body, _ := json.Marshal(fileName)
	req, err := http.NewRequest("POST", "/create-file", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&Server{}).createFileHandler)
	handler.ServeHTTP(rr, req)

	// Save content to the file
	content := []byte("Hello, World!")
	req, err = http.NewRequest("POST", "/files/"+fileName, bytes.NewBuffer(content))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc((&Server{}).fileHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("FileHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Retrieve the file content
	req, err = http.NewRequest("GET", "/files/"+fileName, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc((&Server{}).fileHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("FileHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response FileResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Content != string(content) {
		t.Errorf("FileHandler returned unexpected content: got %v want %v", response.Content, string(content))
	}
}

func TestRunCodeHandler(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to run code
	req, err := http.NewRequest("GET", "/runcode?lang=go", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&Server{}).Runcode)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("RunCodeHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
