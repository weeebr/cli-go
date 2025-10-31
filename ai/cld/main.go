package main

// DESCRIPTION: Claude

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

	// Detect input mode
	inputMode := ai.DetectInputMode()

	var message string
	var err error

	switch inputMode {
	case ai.InputArgs:
		message = ai.GetArgs()
	case ai.InputStdin:
		message, err = ai.ReadStdin()
		ai.ExitIf(err, "failed to read stdin")
	}

	if message == "" {
		ai.LogError("no message provided")
		os.Exit(1)
	}

	// Create Claude client
	// Get Anthropic key (READ-ONLY)
	apiKey, err := config.GetKey("anthropic")
	ai.ExitIf(err, "failed to get Anthropic API key")

	// Create Claude client with API key
	client := ai.NewClaudeClient(ai.ClaudeSonnet, apiKey)

	// Send message with timing
	response, responseInfo, err := ai.SendMessageWithTiming(client, message)
	ai.ExitIf(err, "failed to send message to Claude")

	// Format output based on --json flag
	if toolConfig.JSON {
		// JSON output when --json flag is provided
		jsonData := map[string]string{
			"response": response,
			"model":    responseInfo.Model,
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

	flags.ReorderAndParse()

	return toolConfig
}
