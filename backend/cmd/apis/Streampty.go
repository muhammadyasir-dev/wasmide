package apis

import (
	"bytes"
	"fmt"

	"io"
	//"muhammadyasir-dev/cmd/utils"

	"net/http"
	"os/exec"
	"strings"
)

func isContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "ps", "--filter", "name="+containerName, "--filter", "status=running", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}
func startContainer(containerName string) error {
	cmd := exec.Command("docker", "start", containerName)
	return cmd.Run()
}

func containerExists(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "--filter", "name="+containerName, "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func executeCommand(projectName, command string) (string, error) {
	if projectName == "" {
		return "", fmt.Errorf("project name cannot be empty")
	}

	containerName := fmt.Sprintf("container-%s", projectName)
	var out bytes.Buffer

	// Check container state and manage lifecycle
	if containerExists(containerName) {
		if !isContainerRunning(containerName) {
			fmt.Printf("Starting existing container: %s\n", containerName)
			if err := startContainer(containerName); err != nil {
				return "", fmt.Errorf("failed to start existing container: %v", err)
			}
		} else {
			fmt.Printf("Using running container: %s\n", containerName)
		}
	} else {
		// Create and start new container
		createCmd := exec.Command("docker", "run", "--name", containerName, "-d", "debian:buster-slim", "sleep", "infinity")
		if err := createCmd.Run(); err != nil {
			return "", fmt.Errorf("failed to create container: %v", err)
		}
		fmt.Printf("Created new container: %s\n", containerName)
	}

	// Execute the command in the container
	cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("command execution failed: %v\nOutput: %s", err, out.String())
	}

	return out.String(), nil
}

func Streampty(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, text")              // Allowed headers
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	projectName := r.URL.Query().Get("project")
	if projectName == "" {
		http.Error(w, "Project name is required in query parameters", http.StatusBadRequest)
		return
	}

	command, err := io.ReadAll(r.Body)

	//put things in a queue
	//utils.Rabbitmqueproducer()

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	commandStr := strings.TrimSpace(string(command))
	if commandStr == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	output, err := executeCommand(projectName, commandStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing command: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(output))
}
