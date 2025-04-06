package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFileDir = "./test_files"

func TestMain(m *testing.M) {
	// Setup test environment
	os.Mkdir(testFileDir, 0755)
	defer os.RemoveAll(testFileDir)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateFileHandler(t *testing.T) {
	// Setup test server
	server := &Server{logger: log.New(io.Discard, "", 0)}
	fileDir = testFileDir

	t.Run("Valid file creation", func(t *testing.T) {
		fileName := "testfile.txt"
		body, _ := json.Marshal(fileName)
		req := httptest.NewRequest("POST", "/create-file", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		server.createFileHandler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response FileResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.FileExists(t, filepath.Join(testFileDir, fileName))
	})

	t.Run("Invalid filename", func(t *testing.T) {
		body, _ := json.Marshal("invalid/file.txt")
		req := httptest.NewRequest("POST", "/create-file", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		server.createFileHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestFileHandler(t *testing.T) {
	server := &Server{logger: log.New(io.Discard, "", 0)}
	fileDir = testFileDir

	// Create test file
	testContent := "test content"
	fileName := "testfile.txt"
	os.WriteFile(filepath.Join(testFileDir, fileName), []byte(testContent), 0644)

	t.Run("GET existing file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/files/"+fileName, nil)
		w := httptest.NewRecorder()

		server.fileHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response FileResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, testContent, response.Content)
	})

	t.Run("POST update file", func(t *testing.T) {
		newContent := "new content"
		req := httptest.NewRequest("POST", "/files/"+fileName, bytes.NewBufferString(newContent))
		w := httptest.NewRecorder()

		server.fileHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify file content
		content, _ := os.ReadFile(filepath.Join(testFileDir, fileName))
		assert.Equal(t, newContent, string(content))
	})
}

func TestListFilesHandler(t *testing.T) {
	server := &Server{logger: log.New(io.Discard, "", 0)}
	fileDir = testFileDir

	// Create test files
	os.WriteFile(filepath.Join(testFileDir, "file1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testFileDir, "file2.txt"), []byte(""), 0644)

	req := httptest.NewRequest("GET", "/list-files", nil)
	w := httptest.NewRecorder()

	server.listFilesHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var files []string
	json.Unmarshal(w.Body.Bytes(), &files)
	assert.Contains(t, files, "file1.txt")
	assert.Contains(t, files, "file2.txt")
}

func TestRunCodeHandler(t *testing.T) {
	server := &Server{logger: log.New(io.Discard, "", 0)}

	t.Run("Valid code execution", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/runcode?lang=go", nil)
		w := httptest.NewRecorder()

		server.Runcode(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing language parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/runcode", nil)
		w := httptest.NewRecorder()

		server.Runcode(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCorsMiddleware(t *testing.T) {
	server := &Server{logger: log.New(io.Discard, "", 0)}

	req := httptest.NewRequest("OPTIONS", "/any-endpoint", nil)
	w := httptest.NewRecorder()

	handler := server.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {})
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
