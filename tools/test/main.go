package main

// DESCRIPTION: run e2e tests from config.yml

import (
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	// Get filter from first arg
	var filter string
	if len(os.Args) > 1 {
		filter = os.Args[1]
	}

	// Load config
	configData, err := config.LoadConfig()
	ai.ExitIf(err, "failed to load configuration")

	// Get project root for bin directory
	projectRoot, err := os.Getwd()
	ai.ExitIf(err, "failed to get current directory")
	binDir := filepath.Join(projectRoot, "bin")

	// Check if bin directory exists
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		io.LogError("bin directory not found. Run 'go run build.go' first")
		os.Exit(1)
	}

	// Run tests
	executed := 0
	failed := 0

	io.LogInfo("Running tests from config.yml")
	if filter != "" {
		io.LogInfo("Filtering for commands starting with: %s", filter)
	}

	for _, cmd := range configData.Tests {
		// Apply filter if specified
		if filter != "" && !strings.HasPrefix(cmd, filter) {
			continue
		}

		io.LogInfo("Executing: %s", cmd)

		// Parse command into binary and args
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			io.LogWarning("Empty command, skipping")
			continue
		}

		binary := parts[0]
		args := parts[1:]

		// Build full path to binary
		binaryPath := filepath.Join(binDir, binary)

		// Check if binary exists
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			io.LogError("Binary not found: %s", binaryPath)
			failed++
			continue
		}

		// Execute command
		cmdObj := exec.Command(binaryPath, args...)
		cmdObj.Stdout = os.Stdout
		cmdObj.Stderr = os.Stderr

		if err := cmdObj.Run(); err != nil {
			io.LogError("Command failed: %s", cmd)
			failed++
		} else {
			io.LogSuccess("Command succeeded: %s", cmd)
		}

		executed++
		fmt.Println() // Add spacing between commands
	}

	// Summary
	fmt.Println("==========================================")
	io.LogInfo("Test Summary:")
	fmt.Printf("  Executed: %d\n", executed)
	fmt.Printf("  Failed: %d\n", failed)
	fmt.Printf("  Total tests: %d\n", len(configData.Tests))

	if failed > 0 {
		os.Exit(failed)
	}
}
