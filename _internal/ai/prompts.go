package ai

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cli-go/_internal/config"
)

// PromptClient handles prompt file discovery and selection
type PromptClient struct {
	baseDir string
}

// NewPromptClient creates a new prompt client
func NewPromptClient() *PromptClient {
	// Load config to get prompts directory
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fallback to default if config can't be loaded
		homeDir, err := os.UserHomeDir()
		if err != nil {
			baseDir := filepath.Join(homeDir, ".prompts")
			return &PromptClient{baseDir: baseDir}
		}
		baseDir := filepath.Join(homeDir, "dev", "_private", "prompts-manager", "prompts")
		return &PromptClient{baseDir: baseDir}
	}

	// Use config value, expanding ~ if present
	baseDir := cfg.Prompts.BaseDir
	if strings.HasPrefix(baseDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			baseDir = strings.Replace(baseDir, "~", homeDir, 1)
		}
	}

	return &PromptClient{
		baseDir: baseDir,
	}
}

// SelectPrompt interactively selects a prompt using fzf
func (p *PromptClient) SelectPrompt() (string, error) {
	// Find all .md files in the prompts directory
	files, err := p.findPromptFiles()
	if err != nil {
		return "", fmt.Errorf("failed to find prompt files: %v", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no prompt files found")
	}

	// Format files for fzf display
	formattedFiles := p.formatFilesForFzf(files)

	// Use fzf to select a file
	selected, err := p.runFzf(formattedFiles)
	if err != nil {
		return "", fmt.Errorf("fzf selection failed: %v", err)
	}

	// Convert back to actual file path
	actualFile := p.convertFzfToActualPath(selected)
	return actualFile, nil
}

// LoadPrompt loads the content of a prompt file
func (p *PromptClient) LoadPrompt(promptFile string) (string, error) {
	content, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %v", err)
	}

	return string(content), nil
}

// findPromptFiles finds all .md files in the prompts directory
func (p *PromptClient) findPromptFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(p.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// formatFilesForFzf formats file paths for fzf display
func (p *PromptClient) formatFilesForFzf(files []string) []string {
	var formatted []string

	for _, file := range files {
		// Remove base directory
		relative := strings.TrimPrefix(file, p.baseDir+"/")

		// Remove .md extension
		relative = strings.TrimSuffix(relative, ".md")

		// Replace - with spaces
		relative = strings.ReplaceAll(relative, "-", " ")

		// Replace / with →
		relative = strings.ReplaceAll(relative, "/", " → ")

		formatted = append(formatted, relative)
	}

	return formatted
}

// runFzf runs fzf with the given options
func (p *PromptClient) runFzf(options []string) (string, error) {
	cmd := exec.Command("fzf", "--prompt", "Select prompt → ")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// isInteractive checks if we're in an interactive environment
func isInteractive() bool {
	// Check if we have a TTY
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	// Check if stdin is a character device (TTY)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// convertFzfToActualPath converts fzf selection back to actual file path
func (p *PromptClient) convertFzfToActualPath(fzfSelection string) string {
	// Reverse the formatting
	actual := strings.ReplaceAll(fzfSelection, " → ", "/")
	actual = strings.ReplaceAll(actual, " ", "-")
	actual = actual + ".md"

	return filepath.Join(p.baseDir, actual)
}
