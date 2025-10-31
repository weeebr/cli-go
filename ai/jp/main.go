package main

// DESCRIPTION: ChatGPT w/ prompts

import (
	"flag"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"os"
	"path/filepath"
	"strings"
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

	// Create prompt client
	promptClient := ai.NewPromptClient()

	// Select prompt (use --prompt flag if provided, test mode, or fzf)
	var promptFile string
	var err error

	if toolConfig.Prompt != "" {
		// Use provided prompt file
		promptFile = toolConfig.Prompt
	} else if toolConfig.Test {
		// Test mode: use specific test prompt
		config, err := config.LoadConfig()
		ai.ExitIf(err, "failed to load config")
		promptFile = filepath.Join(config.Prompts.BaseDir, "tools", "translate.md")
	} else {
		// Use fzf to select prompt
		promptFile, err = promptClient.SelectPrompt()
		ai.ExitIf(err, "failed to select prompt")
	}

	// Load prompt content
	promptContent, err := promptClient.LoadPrompt(promptFile)
	ai.ExitIf(err, "failed to load prompt")

	// Get additional message if provided
	var additionalMessage string
	if len(os.Args) > 1 {
		// Filter out --json flag
		var filteredArgs []string
		for _, arg := range os.Args[1:] {
			if arg != "--json" {
				filteredArgs = append(filteredArgs, arg)
			}
		}
		additionalMessage = strings.Join(filteredArgs, " ")
	}

	// If no additional message, get it interactively
	if additionalMessage == "" {
		interactive := io.NewInteractiveInput()
		additionalMessage, err = interactive.GetInput("Enter your message (or Ctrl+D to start conversation)...")
		ai.ExitIf(err, "failed to get input")
	}

	// Add JSON format instruction if --json flag is used
	if toolConfig.JSON {
		promptContent = promptContent + "\n\nPlease respond in valid JSON format."
	}

	// Get credential store (READ-ONLY)

	// Get OpenAI key from store (READ-ONLY)
	apiKey, err := config.GetKey("openai")
	ai.ExitIf(err, "failed to get OpenAI API key")

	// Create ChatGPT client with API key
	format := "text"
	if toolConfig.JSON {
		format = "json"
	}
	client := ai.NewChatGPTClient(format, apiKey)

	// Send message with role file and track timing
	start := time.Now()
	response, err := client.SendMessageWithRoleFile(promptFile, additionalMessage)
	duration := time.Since(start)
	responseInfo := ai.ResponseInfo{
		Duration: duration,
		Model:    client.GetModel(),
	}

	ai.ExitIf(err, "failed to send message to ChatGPT")

	// Format output based on --json flag
	if toolConfig.JSON {
		// JSON output when --json flag is provided
		// Extract prompt name from file path
		var promptName string
		// Extract prompt name by removing base directory and .md extension
		config, err := config.LoadConfig()
		if err == nil {
			baseDir := config.Prompts.BaseDir
			if strings.HasPrefix(promptFile, baseDir) {
				promptName = strings.TrimSuffix(strings.TrimPrefix(promptFile, baseDir+"/"), ".md")
			}
		}
		if promptName == "" {
			// Fallback: extract from filename
			promptName = strings.TrimSuffix(filepath.Base(promptFile), ".md")
		}

		jsonData := map[string]interface{}{
			"response":    response,
			"prompt_used": promptName,
			"prompt_file": promptFile,
			"format":      format,
			"model":       responseInfo.Model,
		}

		// Use direct output
		io.DirectOutput(jsonData, toolConfig.Clip, toolConfig.File, toolConfig.JSON)
	} else {
		// Default: markdown output formatted with glamour and response info
		responseInfoStr := ai.FormatResponseInfo(responseInfo)
		io.FormatTerminalOutputWithResponseInfo(response, responseInfoStr)
	}
}

func parseFlags() ToolConfig {
	toolConfig := ToolConfig{}

	flag.BoolVar(&toolConfig.Clip, "clip", false, "Copy to clipboard")
	flag.StringVar(&toolConfig.File, "file", "", "Write to file")
	flag.BoolVar(&toolConfig.JSON, "json", false, "Output in JSON format")
	flag.StringVar(&toolConfig.Prompt, "prompt", "", "Path to prompt file")
	flag.BoolVar(&toolConfig.Test, "test", false, "Test mode - use translate.md prompt")

	flags.ReorderAndParse()

	// Check for --json flag in args
	for _, arg := range os.Args {
		if arg == "--json" {
			toolConfig.JSON = true
			break
		}
	}

	return toolConfig
}
