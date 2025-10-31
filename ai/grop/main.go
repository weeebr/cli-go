package main

// DESCRIPTION: Grok w/ prompts

import (
	"flag"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/io"
	"os"
	"path/filepath"
	"strings"
)

type ToolConfig struct {
	Clip   bool
	File   string
	JSON   bool
	Prompt string
	Test   bool
}

func main() {
	toolConfig := parseFlags()

	// Check for help command
	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		io.LogInfo("grop - Grok with prompt selection")
		io.LogInfo("Selects a prompt using fzf and uses it with Grok")
		io.LogInfo("Usage: grop 'additional context' | grop")
		io.LogInfo("Output: {\"response\": \"Extracted content\", \"prompt_used\": \"prompt-name\", \"model\": \"grok-code-fast-1\"}")
		return
	}

	// Select prompt (use --prompt flag if provided, test mode, or fzf)
	var promptFile string
	var err error

	if toolConfig.Prompt != "" {
		// Use provided prompt file for testing
		promptFile = toolConfig.Prompt
	} else if toolConfig.Test {
		// Test mode: use specific test prompt
		config, err := config.LoadConfig()
		ai.ExitIf(err, "failed to load config")
		promptFile = filepath.Join(config.Prompts.BaseDir, "tools", "translate.md")
	} else {
		// Use fzf to select prompt
		promptClient := ai.NewPromptClient()
		promptFile, err = promptClient.SelectPrompt()
		ai.ExitIf(err, "failed to select prompt")
	}

	// Get additional message if provided (filter out flags)
	var additionalMessage string
	var filteredArgs []string
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "--") && !strings.HasPrefix(arg, "-") {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	if len(filteredArgs) > 0 {
		additionalMessage = strings.Join(filteredArgs, " ")
	}

	// If no additional message, use a default message
	if additionalMessage == "" {
		additionalMessage = "Please help me with this task."
	}

	// Create Grok client (no API key needed for CLI wrapper)
	client := ai.NewGrokClient()

	// Send message with prompt
	response, err := client.SendMessageWithPrompt(promptFile, additionalMessage)
	ai.ExitIf(err, "failed to send message to Grok")

	// Detect files in response
	files := client.DetectFiles(response)

	// Cleanup files if any were detected
	if len(files) > 0 {
		client.CleanupFiles(files)
	}

	// Format output based on --json flag
	if toolConfig.JSON {
		// JSON output when --json flag is provided
		jsonData := map[string]string{
			"response":    response,
			"prompt_used": promptFile,
			"model":       client.GetModel(),
		}

		// Use direct output
		io.DirectOutput(jsonData, toolConfig.Clip, toolConfig.File, toolConfig.JSON)
	} else {
		// Default: markdown output formatted with glamour
		io.FormatTerminalOutput(response)
	}

	// Cleanup .grok directory
	defer client.Cleanup()
}

func parseFlags() ToolConfig {
	toolConfig := ToolConfig{}

	// Add basic flag support for consistency
	flag.BoolVar(&toolConfig.Clip, "clip", false, "Copy to clipboard")
	flag.StringVar(&toolConfig.File, "file", "", "Write to file")
	flag.BoolVar(&toolConfig.JSON, "json", false, "Output in JSON format")
	flag.StringVar(&toolConfig.Prompt, "prompt", "", "Prompt file path (for testing, bypasses fzf)")
	flag.BoolVar(&toolConfig.Test, "test", false, "Test mode - use translate.md prompt and default message")

	flag.Parse()

	return toolConfig
}
