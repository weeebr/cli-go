package main

// DESCRIPTION: ChatGPT (JSON)

import (
	"flag"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"os"
	"time"
)

type ToolConfig struct {
	Clip   bool
	File   string
	JSON   bool
	Prompt string
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

	// Get OpenAI key from store (READ-ONLY)
	apiKey, err := config.GetKey("openai")
	ai.ExitIf(err, "failed to get OpenAI API key")

	// Create ChatGPT client with API key
	client := ai.NewChatGPTClient("json", apiKey)

	// Send message with or without prompt file and track response time
	var response string
	var responseInfo ai.ResponseInfo

	if toolConfig.Prompt != "" {
		// Use prompt file - need to track timing manually
		start := time.Now()
		response, err = client.SendMessageWithRoleFile(toolConfig.Prompt, message)
		duration := time.Since(start)
		responseInfo = ai.ResponseInfo{
			Duration: duration,
			Model:    client.GetModel(),
		}
	} else {
		// Send message directly with timing
		response, responseInfo, err = ai.SendMessageWithTiming(client, message)
	}

	ai.ExitIf(err, "failed to send message to ChatGPT")

	// Format output based on --json flag
	if toolConfig.JSON {
		// JSON output when --json flag is provided
		jsonData := map[string]interface{}{
			"data":   response,
			"model":  responseInfo.Model,
			"format": "json",
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

	flags.ReorderAndParse()

	return toolConfig
}
