package main

// DESCRIPTION: ChatGPT

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
	case ai.InputInteractive:
		// Use input.md pattern for interactive input
		interactive := io.NewInteractiveInput()
		message, err = interactive.GetInput("Enter your message (or Ctrl+D to start conversation)...")
		ai.ExitIf(err, "failed to get input")
	}

	if message == "" {
		ai.LogError("no message provided")
		os.Exit(1)
	}

	// Get credential store (READ-ONLY)

	// Get OpenAI key from store (READ-ONLY)
	apiKey, err := config.GetKey("openai")
	ai.ExitIf(err, "failed to get OpenAI API key")

	// Create ChatGPT client with API key
	client := ai.NewChatGPTClient("text", apiKey)

	// Send message
	response, err := client.SendMessage(message)
	ai.ExitIf(err, "failed to send message to ChatGPT")

	// Format output based on --json flag
	if toolConfig.JSON {
		// JSON output when --json flag is provided
		jsonData := map[string]string{
			"response": response,
			"model":    client.GetModel(),
			"thread":   client.GetThread(),
		}

		// Use direct output
		io.DirectOutput(jsonData, toolConfig.Clip, toolConfig.File, toolConfig.JSON)
	} else {
		// Default: markdown output formatted with glamour
		io.FormatTerminalOutput(response)
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
