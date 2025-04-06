package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xxx/runnerservice"
)

// Configuration constants
const (
	fileDir      = "./files" // Directory to store files
	maxFileSizes = 10 << 20  // 10 MB maximum file size
	serverPort   = ":8082"   // Server port
)

// FileChange represents a change made to a file
type FileChange struct {
	Type      string    `json:"type"`      // Type of change (added/removed)
	Content   string    `json:"content"`   // Changed content
	Timestamp time.Time `json:"timestamp"` // When the change occurred
}

// FileResponse represents the response structure for file operations
type FileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Content string `json:"content,omitempty"`
}

// Server represents our HTTP server and its dependencies
type Server struct {
	logger *log.Logger
}

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[FileEditor] ", log.LstdFlags|log.Lshortfile)

	// Create new server instance
	server := &Server{
		logger: logger,
	}

	// Ensure the files directory exists
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		logger.Fatalf("Failed to create files directory: %v", err)
	}

	// Initialize routes
	mux := http.NewServeMux()
	mux.HandleFunc("/files/", server.corsMiddleware(server.fileHandler))
	mux.HandleFunc("/create-file", server.corsMiddleware(server.createFileHandler))
	mux.HandleFunc("/list-files", server.corsMiddleware(server.listFilesHandler))
	mux.HandleFunc("/runcode", server.corsMiddleware(server.Runcode))
	// Configure server
	srv := &http.Server{
		Addr:         serverPort,
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  6 * time.Second,
	}

	// Start server
	logger.Printf("Starting server on port %s...", serverPort)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

// corsMiddleware handles CORS headers and preflight requests
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// fileHandler handles operations on individual files
func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)

	// Validate filename
	if !isValidFilename(fileName) {
		s.jsonResponse(w, http.StatusBadRequest, FileResponse{
			Success: false,
			Message: "Invalid filename",
		})
		return
	}

	filePath := filepath.Join(fileDir, fileName)

	switch r.Method {
	case http.MethodGet:
		s.handleGetFile(w, filePath)
	case http.MethodPost:
		s.handleSaveFile(w, r, filePath)
	default:
		s.jsonResponse(w, http.StatusMethodNotAllowed, FileResponse{
			Success: false,
			Message: "Method not allowed",
		})
	}
}

// handleGetFile handles retrieving file content
func (s *Server) handleGetFile(w http.ResponseWriter, filePath string) {
	content, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		s.jsonResponse(w, http.StatusNotFound, FileResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}
	if err != nil {
		s.logger.Printf("Error reading file %s: %v", filePath, err)
		s.jsonResponse(w, http.StatusInternalServerError, FileResponse{
			Success: false,
			Message: "Error reading file",
		})
		return
	}

	s.jsonResponse(w, http.StatusOK, FileResponse{
		Success: true,
		Content: string(content),
	})
}

// handleSaveFile handles saving file content
func (s *Server) handleSaveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	// Limit request size
	r.Body = http.MaxBytesReader(w, r.Body, 10048)

	content, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("Error reading request body: %v", err)
		s.jsonResponse(w, http.StatusBadRequest, FileResponse{
			Success: false,
			Message: "Error reading request body",
		})
		return
	}

	err = os.WriteFile(filePath, content, 0644)

	if err != nil {
		s.logger.Printf("Error writing file %s: %v", filePath, err)
		s.jsonResponse(w, http.StatusInternalServerError, FileResponse{
			Success: false,
			Message: "Error writing file",
		})
		return
	}

	s.jsonResponse(w, http.StatusOK, FileResponse{
		Success: true,
		Message: "File saved successfully",
	})
}

// createFileHandler handles creating new files
func (s *Server) createFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.jsonResponse(w, http.StatusMethodNotAllowed, FileResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var fileName string
	if err := json.NewDecoder(r.Body).Decode(&fileName); err != nil {
		s.jsonResponse(w, http.StatusBadRequest, FileResponse{
			Success: false,
			Message: "Invalid file name",
		})
		return
	}

	if !isValidFilename(fileName) {
		s.jsonResponse(w, http.StatusBadRequest, FileResponse{
			Success: false,
			Message: "Invalid filename",
		})
		return
	}

	filePath := filepath.Join(fileDir, fileName)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		s.jsonResponse(w, http.StatusConflict, FileResponse{
			Success: false,
			Message: "File already exists",
		})
		return
	}

	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		s.logger.Printf("Error creating file %s: %v", filePath, err)
		s.jsonResponse(w, http.StatusInternalServerError, FileResponse{
			Success: false,
			Message: "Error creating file",
		})
		return
	}

	s.jsonResponse(w, http.StatusCreated, FileResponse{
		Success: true,
		Message: "File created successfully",
	})
}

// listFilesHandler handles listing all files in the directory
func (s *Server) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.jsonResponse(w, http.StatusMethodNotAllowed, FileResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	files, err := os.ReadDir(fileDir)
	if err != nil {
		s.logger.Printf("Error reading directory: %v", err)
		s.jsonResponse(w, http.StatusInternalServerError, FileResponse{
			Success: false,
			Message: "Error reading directory",
		})
		return
	}

	fileList := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)
}
func (s *Server) Runcode(w http.ResponseWriter, r *http.Request) {
	//first run the assembly files
	programminglang := r.URL.Query().Get("lang")
	runnerservice.Execwasm(programminglang)
	w.WriteHeader(http.StatusOK)
}

// jsonResponse sends a JSON response with the given status code and data
func (s *Server) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Printf("Error encoding JSON response: %v", err)
	}
}

// isValidFilename checks if the filename is valid and safe
func isValidFilename(filename string) bool {
	// Check for empty filename
	if filename == "" {
		return false
	}

	// Check for directory traversal attempts
	if strings.Contains(filename, "..") {
		return false
	}

	// Check for invalid characters
	return !strings.ContainsAny(filename, "/\\?%*:|\"<>")
}
