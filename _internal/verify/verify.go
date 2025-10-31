package verify

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// VerificationResult represents the result of a verification check
type VerificationResult struct {
	Tool      string   `json:"tool"`
	Status    string   `json:"status"`
	Error     string   `json:"error,omitempty"`
	Duration  string   `json:"duration"`
	Checks    []string `json:"checks"`
	Timestamp string   `json:"timestamp"`
}

// VerificationSummary represents the overall verification summary
type VerificationSummary struct {
	Total     int                  `json:"total"`
	Passed    int                  `json:"passed"`
	Failed    int                  `json:"failed"`
	Results   []VerificationResult `json:"results"`
	Duration  string               `json:"duration"`
	Timestamp string               `json:"timestamp"`
}

// FindAllTools finds all tools in the project
func FindAllTools(projectRoot string) []string {
	var tools []string

	// Check git/ directory
	gitDir := filepath.Join(projectRoot, "git")
	if entries, err := os.ReadDir(gitDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				tools = append(tools, entry.Name())
			}
		}
	}

	// Check ai/ directory
	aiDir := filepath.Join(projectRoot, "ai")
	if entries, err := os.ReadDir(aiDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				tools = append(tools, entry.Name())
			}
		}
	}

	// Check tools/ directory
	toolsDir := filepath.Join(projectRoot, "tools")
	if entries, err := os.ReadDir(toolsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				tools = append(tools, entry.Name())
			}
		}
	}

	return tools
}

// RunVerification runs verification on all tools
func RunVerification(tools []string, projectRoot string) []VerificationResult {
	var results []VerificationResult

	for _, toolName := range tools {
		startTime := time.Now()
		result := VerificationResult{
			Tool:      toolName,
			Timestamp: startTime.Format(time.RFC3339),
		}

		// Find tool directory
		toolDir := findToolDirectory(toolName, projectRoot)
		if toolDir == "" {
			result.Status = "error"
			result.Error = "Tool directory not found"
			result.Duration = time.Since(startTime).String()
			results = append(results, result)
			continue
		}

		// Run verification checks
		var checks []string
		var hasErrors bool

		// Go vet
		if err := runGoVet(toolDir); err != nil {
			checks = append(checks, fmt.Sprintf("go vet: %v", err))
			hasErrors = true
		} else {
			checks = append(checks, "go vet: passed")
		}

		// Staticcheck
		if err := runStaticcheck(toolDir); err != nil {
			checks = append(checks, fmt.Sprintf("staticcheck: %v", err))
			hasErrors = true
		} else {
			checks = append(checks, "staticcheck: passed")
		}

		// Self test
		if err := runSelfTest(toolName, projectRoot); err != nil {
			checks = append(checks, fmt.Sprintf("self test: %v", err))
			hasErrors = true
		} else {
			checks = append(checks, "self test: passed")
		}

		if hasErrors {
			result.Status = "failed"
		} else {
			result.Status = "passed"
		}

		result.Checks = checks
		result.Duration = time.Since(startTime).String()
		results = append(results, result)
	}

	return results
}

// CreateSummary creates a verification summary from results
func CreateSummary(results []VerificationResult, startTime time.Time) VerificationSummary {
	summary := VerificationSummary{
		Total:     len(results),
		Results:   results,
		Duration:  time.Since(startTime).String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	for _, result := range results {
		if result.Status == "passed" {
			summary.Passed++
		} else {
			summary.Failed++
		}
	}

	return summary
}

func findToolDirectory(toolName, projectRoot string) string {
	// Check git/ directory
	gitDir := filepath.Join(projectRoot, "git", toolName)
	if _, err := os.Stat(gitDir); err == nil {
		return gitDir
	}

	// Check ai/ directory
	aiDir := filepath.Join(projectRoot, "ai", toolName)
	if _, err := os.Stat(aiDir); err == nil {
		return aiDir
	}

	// Check tools/ directory
	toolsDir := filepath.Join(projectRoot, "tools", toolName)
	if _, err := os.Stat(toolsDir); err == nil {
		return toolsDir
	}

	return ""
}

func runGoVet(toolDir string) error {
	cmd := exec.Command("go", "vet", "./"+filepath.Base(toolDir))
	cmd.Dir = filepath.Dir(toolDir)
	return cmd.Run()
}

func runStaticcheck(toolDir string) error {
	cmd := exec.Command("staticcheck", "./"+filepath.Base(toolDir))
	cmd.Dir = filepath.Dir(toolDir)
	return cmd.Run()
}

func runSelfTest(toolName, projectRoot string) error {
	// Build the tool first
	buildCmd := exec.Command("go", "build", "-o", "/tmp/verify-"+toolName, "./"+filepath.Join("tools", toolName))
	buildCmd.Dir = projectRoot
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	// Run the tool with --help to test basic functionality
	testCmd := exec.Command("/tmp/verify-"+toolName, "--help")
	testCmd.Dir = projectRoot
	if err := testCmd.Run(); err != nil {
		return fmt.Errorf("self test failed: %v", err)
	}

	// Clean up
	os.Remove("/tmp/verify-" + toolName)
	return nil
}
