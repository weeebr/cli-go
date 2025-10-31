package main

// DESCRIPTION: open prompts in Cursor

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/io"
	"os"
	"os/exec"
)

type PromptsResult struct {
	Action  string `json:"action"`
	Path    string `json:"path,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
}

func main() {

	var (
		clip   = flag.Bool("clip", false, "Copy to clipboard")
		file   = flag.String("file", "", "Write to file")
		editor = flag.String("editor", "cursor", "Editor to use (cursor, code, vim, etc.)")
		json   = flag.Bool("json", false, "Output in JSON format")
	)
	flag.Parse()

	// Load config to get prompts directory
	config, err := config.LoadConfig()
	ai.ExitIf(err, "failed to load config")
	promptsDir := config.Prompts.BaseDir

	// Open prompts in editor
	var cmd *exec.Cmd
	switch *editor {
	case "cursor":
		cmd = exec.Command("cursor", promptsDir)
	case "code":
		cmd = exec.Command("code", promptsDir)
	case "vim":
		cmd = exec.Command("vim", promptsDir)
	case "nvim":
		cmd = exec.Command("nvim", promptsDir)
	default:
		cmd = exec.Command(*editor, promptsDir)
	}

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result := PromptsResult{
			Action:  "open_prompts",
			Path:    promptsDir,
			Success: false,
			Error:   fmt.Sprintf("failed to open prompts: %v", err),
			Message: fmt.Sprintf("Failed to open prompts in %s", *editor),
		}
		// Use direct output
		if *json {
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			// Default: formatted output
			io.FormatTerminalOutput(result.Message)
		}
		os.Exit(1)
	}

	result := PromptsResult{
		Action:  "open_prompts",
		Path:    promptsDir,
		Success: true,
		Message: fmt.Sprintf("Opened prompts directory in %s", *editor),
	}

	// Use direct output
	if *json {
		io.DirectOutput(result, *clip, *file, *json)
	} else {
		// Default: formatted output
		io.FormatTerminalOutput(result.Message)
	}
}
