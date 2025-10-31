package main

// DESCRIPTION: Grok

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"os"
	"path/filepath"
	"time"
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

	// Handle prompt selection
	var promptFile string
	if toolConfig.Prompt != "" {
		// Use provided prompt path
		promptFile = toolConfig.Prompt
	} else if toolConfig.Test {
		// Test mode: use specific test prompt
		config, err := config.LoadConfig()
		ai.ExitIf(err, "failed to load config")
		promptFile = filepath.Join(config.Prompts.BaseDir, "tools", "translate.md")
	} else {
		// Interactive prompt selection
		promptClient := ai.NewPromptClient()
		selectedPrompt, err := promptClient.SelectPrompt()
		ai.ExitIf(err, "failed to select prompt")
		promptFile = selectedPrompt
	}

	// Detect input mode
	inputMode := ai.DetectInputMode()
	fmt.Fprintf(os.Stderr, "DEBUG: inputMode = %s\n", inputMode)

	var message string
	var err error

	switch inputMode {
	case ai.InputArgs:
		message = ai.GetArgs()
		fmt.Fprintf(os.Stderr, "DEBUG: message from args = %s\n", message)
	case ai.InputStdin:
		message, err = ai.ReadStdin()
		ai.ExitIf(err, "failed to read stdin")
		fmt.Fprintf(os.Stderr, "DEBUG: message from stdin = %s\n", message)
	case ai.InputInteractive:
		// Use input.md pattern for interactive input
		interactive := io.NewInteractiveInput()
		message, err = interactive.GetInput("Enter your message (or Ctrl+D to start conversation)...")
		ai.ExitIf(err, "failed to get input")
		fmt.Fprintf(os.Stderr, "DEBUG: message from interactive = %s\n", message)
	}

	if message == "" {
		ai.LogError("no message provided")
		os.Exit(1)
	}

	// Create Grok client (no API key needed for CLI wrapper)
	client := ai.NewGrokClient()

	// Send message with prompt and track timing
	start := time.Now()
	response, err := client.SendMessageWithPrompt(promptFile, message)
	duration := time.Since(start)
	responseInfo := ai.ResponseInfo{
		Duration: duration,
		Model:    client.GetModel(),
	}

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
		jsonData := map[string]interface{}{
			"response":       response,
			"model":          responseInfo.Model,
			"files_detected": files,
		}

		// Use direct output
		io.DirectOutput(jsonData, toolConfig.Clip, toolConfig.File, toolConfig.JSON)
	} else {
		// Default: markdown output formatted with glamour and response info
		responseInfoStr := ai.FormatResponseInfo(responseInfo)
		io.FormatTerminalOutputWithResponseInfo(response, responseInfoStr)
	}

	// Cleanup .grok directory
	defer client.Cleanup()
}

func parseFlags() ToolConfig {
	toolConfig := ToolConfig{}

	flag.BoolVar(&toolConfig.Clip, "clip", false, "Copy to clipboard")
	flag.StringVar(&toolConfig.File, "file", "", "Write to file")
	flag.BoolVar(&toolConfig.JSON, "json", false, "Output in JSON format")
	flag.StringVar(&toolConfig.Prompt, "prompt", "", "Path to prompt file (skips interactive selection)")
	flag.BoolVar(&toolConfig.Test, "test", false, "Test mode - use translate.md prompt")

	flags.ReorderAndParse()

	return toolConfig
}
