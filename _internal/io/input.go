package io

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InteractiveInput handles interactive input using input.md pattern
type InteractiveInput struct{}

// NewInteractiveInput creates a new interactive input handler
func NewInteractiveInput() *InteractiveInput {
	return &InteractiveInput{}
}

// GetInput gets input from user using input.md pattern
func (i *InteractiveInput) GetInput(placeholder string) (string, error) {
	// Remove existing input.md
	os.Remove("input.md")

	// Create input.md with placeholder as comment
	file, err := os.Create("input.md")
	if err != nil {
		return "", fmt.Errorf("failed to create input.md: %v", err)
	}

	// Write placeholder as comment
	if placeholder != "" {
		file.WriteString(fmt.Sprintf("<!-- %s -->\n\n", placeholder))
	}
	file.Close()

	// Open in editor
	if err := i.EditFile("input.md"); err != nil {
		return "", fmt.Errorf("failed to edit input.md: %v", err)
	}

	// Read content
	content, err := os.ReadFile("input.md")
	if err != nil {
		return "", fmt.Errorf("failed to read input.md: %v", err)
	}

	// Clean up
	os.Remove("input.md")

	// Remove placeholder comment if it exists
	contentStr := string(content)
	if strings.HasPrefix(contentStr, fmt.Sprintf("<!-- %s -->", placeholder)) {
		lines := strings.Split(contentStr, "\n")
		if len(lines) > 1 {
			contentStr = strings.Join(lines[1:], "\n")
		}
	}

	return strings.TrimSpace(contentStr), nil
}

// EditFile opens a file in the default editor and waits for completion
func (i *InteractiveInput) EditFile(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // fallback editor
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// WaitForUser waits for user to press Enter
func (i *InteractiveInput) WaitForUser(message string) {
	fmt.Print(message)
	var input string
	fmt.Scanln(&input)
}

// LogError logs an error message to stderr
func LogError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", fmt.Sprintf(format, args...))
}

// LogWarning logs a warning message to stderr
func LogWarning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Warning: %s\n", fmt.Sprintf(format, args...))
}

// LogInfo logs an info message to stderr
func LogInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Info: %s\n", fmt.Sprintf(format, args...))
}

// LogSuccess logs a success message to stderr
func LogSuccess(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Success: %s\n", fmt.Sprintf(format, args...))
}
