package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// GrokClient handles xAI Grok CLI integration
type GrokClient struct {
	workDir string // .grok directory path
}

// NewGrokClient creates a new Grok client
func NewGrokClient() *GrokClient {
	return &GrokClient{workDir: ".grok"}
}

// SetupWorkDir creates .grok/ and settings.json
func (g *GrokClient) SetupWorkDir() error {
	if err := os.MkdirAll(g.workDir, 0755); err != nil {
		return fmt.Errorf("failed to create .grok directory: %v", err)
	}

	settings := map[string]string{"model": "grok-code-fast-1"}
	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %v", err)
	}

	settingsFile := filepath.Join(g.workDir, "settings.json")
	if err := os.WriteFile(settingsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %v", err)
	}

	return nil
}

// SendMessage wraps grok CLI
func (g *GrokClient) SendMessage(message string) (string, error) {
	// Setup work directory
	if err := g.SetupWorkDir(); err != nil {
		return "", err
	}

	// Run grok CLI
	cmd := exec.Command("grok", "-p", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("grok CLI failed: %v", err)
	}

	// Extract JSON content using Python
	content, err := g.extractContent(string(output))
	if err != nil {
		return "", err
	}

	return content, nil
}

// SendMessageWithPrompt copies prompt to GROK.md and sends message
func (g *GrokClient) SendMessageWithPrompt(promptFile, message string) (string, error) {
	if err := g.SetupWorkDir(); err != nil {
		return "", err
	}

	// Copy prompt file to .grok/GROK.md (grok CLI expects prompt here)
	grokMdPath := filepath.Join(g.workDir, "GROK.md")
	promptData, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %v", err)
	}

	if err := os.WriteFile(grokMdPath, promptData, 0644); err != nil {
		return "", fmt.Errorf("failed to write GROK.md: %v", err)
	}

	// Run grok CLI with user message (prompt is already in GROK.md)
	cmd := exec.Command("grok", "-p", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("grok CLI failed: %v", err)
	}

	content, err := g.extractContent(string(output))
	if err != nil {
		return "", err
	}

	return content, nil
}

// extractContent uses pure Go to parse JSON output
func (g *GrokClient) extractContent(output string) (string, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to parse as JSON
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}

		// Check if it's an assistant message
		if role, ok := obj["role"].(string); ok && role == "assistant" {
			if content, ok := obj["content"].(string); ok && content != "" {
				return g.cleanContent(content), nil
			}
		}
	}

	return "", fmt.Errorf("no assistant content found in output")
}

// cleanContent removes tool messages and converts escaped newlines
func (g *GrokClient) cleanContent(content string) string {
	// Remove "Using tools to help you..." lines
	lines := strings.Split(content, "\n")
	var cleaned []string
	for _, line := range lines {
		if !strings.Contains(line, "Using tools to help you") {
			cleaned = append(cleaned, line)
		}
	}

	// Join and convert escaped newlines to actual newlines
	result := strings.Join(cleaned, "\n")
	result = strings.ReplaceAll(result, "\\n", "\n")

	return strings.TrimSpace(result)
}

// DetectFiles detects created files in the response
func (g *GrokClient) DetectFiles(response string) []string {
	var files []string

	// Look for file creation patterns
	patterns := []string{
		`created.*\.[a-zA-Z]+`,
		`[a-zA-Z0-9_-]*\.[a-zA-Z0-9]+`,
		`[a-zA-Z0-9_-]*\.md`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(response, -1)
		files = append(files, matches...)
	}

	return files
}

// CleanupFiles removes detected files
func (g *GrokClient) CleanupFiles(files []string) {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
			LogInfo("🗑️  File deleted: %s", file)
		}
	}
}

// Cleanup removes .grok directory
func (g *GrokClient) Cleanup() error {
	return os.RemoveAll(g.workDir)
}

// GetModel returns the model name
func (g *GrokClient) GetModel() string {
	return "grok-code-fast-1"
}
