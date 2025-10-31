package main

// DESCRIPTION: ChatGPT w/ input.md

import (
	"flag"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"os"
)

type ToolConfig struct {
	Clip bool
	File string
	JSON bool
}

func main() {

	toolConfig := parseFlags()

	// Remove existing input.md
	os.Remove("input.md")

	// Create empty input.md
	file, err := os.Create("input.md")
	ai.ExitIf(err, "failed to create input.md")
	file.Close()

	// Open in editor
	interactive := io.NewInteractiveInput()
	io.LogInfo("📝 Edit input.md in your editor... (press Enter when done)")
	interactive.WaitForUser("")

	if err := interactive.EditFile("input.md"); err != nil {
		ai.ExitIf(err, "failed to edit input.md")
	}

	// Check if file has content
	content, err := os.ReadFile("input.md")
	ai.ExitIf(err, "failed to read input.md")

	if len(content) == 0 {
		io.LogError("⚠️  input.md is empty, nothing to process")
		os.Remove("input.md")
		os.Exit(1)
	}

	// Process with ChatGPT
	io.LogInfo("🤖 Processing with ChatGPT...")
	// Get credential store (READ-ONLY)

	// Get OpenAI key from store (READ-ONLY)
	apiKey, err := config.GetKey("openai")
	ai.ExitIf(err, "failed to get OpenAI API key")

	// Create ChatGPT client with API key
	client := ai.NewChatGPTClient("text", apiKey)
	response, responseInfo, err := ai.SendMessageWithTiming(client, string(content))
	ai.ExitIf(err, "failed to send message to ChatGPT")

	// Format output based on --json flag
	if toolConfig.JSON {
		// JSON output when --json flag is provided
		jsonData := map[string]interface{}{
			"processed": true,
			"file":      "input.md",
			"response":  response,
			"model":     responseInfo.Model,
		}

		// Use direct output
		io.DirectOutput(jsonData, toolConfig.Clip, toolConfig.File, toolConfig.JSON)
	} else {
		// Default: markdown output formatted with glamour and response info
		responseInfoStr := ai.FormatResponseInfo(responseInfo)
		io.FormatTerminalOutputWithResponseInfo(response, responseInfoStr)
	}

	// Cleanup
	os.Remove("input.md")
	io.LogInfo("🗑️  Cleaned up input.md")
}

func parseFlags() ToolConfig {
	toolConfig := ToolConfig{}

	flag.BoolVar(&toolConfig.Clip, "clip", false, "Copy to clipboard")
	flag.StringVar(&toolConfig.File, "file", "", "Write to file")
	flag.BoolVar(&toolConfig.JSON, "json", false, "Output in JSON format")

	flags.ReorderAndParse()

	return toolConfig
}
