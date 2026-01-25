package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req RunCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request: %v", err)
		return
	}

	log.Printf("Received code execution request for language: %s", req.Language)

	// Execute the code in a Docker container
	response := executeCodeInDocker(req.Code, req.Language)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// executeCodeInDocker runs the provided code in a sandboxed Docker container
func executeCodeInDocker(code, language string) RunCodeResponse {
	// Configuration
	const (
		timeout     = 10 * time.Second
		memoryLimit = "128m"
		cpuLimit    = "0.5"
	)

	// Only support Python for now
	if language != "python" && language != "python3" {
		return RunCodeResponse{
			Output:   "",
			Error:    fmt.Sprintf("Unsupported language: %s. Only Python is supported.", language),
			ExitCode: 1,
		}
	}

	// Create a temporary file with the code
	// Use shared directory that's mounted in both the core container and accessible to Docker host
	const codeDir = "/tmp/ironsnake-code"
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		log.Printf("Failed to create code directory: %v", err)
		return RunCodeResponse{
			Output:   "",
			Error:    "Internal error: failed to create code directory",
			ExitCode: 1,
		}
	}
	tmpFile, err := os.CreateTemp(codeDir, "code-*.py")
	if err != nil {
		log.Printf("Failed to create temp file: %v", err)
		return RunCodeResponse{
			Output:   "",
			Error:    "Internal error: failed to create temporary file",
			ExitCode: 1,
		}
	}
	defer os.Remove(tmpFile.Name())

	// Write the code to the temp file
	if _, err := tmpFile.WriteString(code); err != nil {
		log.Printf("Failed to write code to temp file: %v", err)
		return RunCodeResponse{
			Output:   "",
			Error:    "Internal error: failed to write code",
			ExitCode: 1,
		}
	}
	tmpFile.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Build docker run command with security restrictions
	args := []string{
		"run",
		"--rm",              // Remove container after execution
		"--network", "none", // No network access
		"--memory", memoryLimit, // Memory limit
		"--cpus", cpuLimit, // CPU limit
		"--pids-limit", "64", // Limit number of processes
		"--read-only",              // Read-only filesystem
		"--tmpfs", "/tmp:size=10m", // Small writable /tmp
		"--security-opt", "no-new-privileges", // Prevent privilege escalation
		"-v", tmpFile.Name() + ":/code.py:ro", // Mount code file read-only
		"python:3.14-slim",   // Python image
		"python", "/code.py", // Run the code
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()

	// Determine exit code
	exitCode := 0
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return RunCodeResponse{
				Output:   stdout.String(),
				Error:    fmt.Sprintf("Execution timed out after %v", timeout),
				ExitCode: 124, // Standard timeout exit code
			}
		}

		// Try to get the actual exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// Docker or other error
			log.Printf("Docker execution error: %v", err)
			errorMsg := stderr.String()
			if errorMsg == "" {
				errorMsg = err.Error()
			}
			// Check for common Docker errors
			if strings.Contains(errorMsg, "Cannot connect to the Docker daemon") ||
				strings.Contains(errorMsg, "docker: not found") {
				return RunCodeResponse{
					Output:   "",
					Error:    "Docker is not available on the server",
					ExitCode: 1,
				}
			}
			return RunCodeResponse{
				Output:   "",
				Error:    fmt.Sprintf("Execution failed: %s", errorMsg),
				ExitCode: 1,
			}
		}
	}

	// Combine output
	output := stdout.String()
	errorOutput := stderr.String()

	return RunCodeResponse{
		Output:   output,
		Error:    errorOutput,
		ExitCode: exitCode,
	}
}
