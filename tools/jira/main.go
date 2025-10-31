package main

// DESCRIPTION: Jira CLI tool

import (
	"flag"
	"cli-go/_internal/flags"
	"cli-go/_internal/jira"
)

func main() {

	toolConfig := parseFlags()

	// Parse arguments
	args := flag.Args()
	if len(args) == 0 {
		// No args - show current user activity (requires config)
		jiraConfig, apiToken, err := jira.LoadJiraConfig()
		jira.HandleError(err, "Failed to load configuration")
		client := jira.NewClient(
			jiraConfig.BaseURL,
			jiraConfig.Email,
			apiToken,
			jiraConfig.DefaultProject,
		)
		flags := &jira.Flags{
			Clip: toolConfig.Clip,
			File: toolConfig.File,
			JSON: toolConfig.JSON,
			Open: toolConfig.Open,
		}
		result := jira.RouteCommand("user", []string{}, client, flags)
		jira.OutputJSON(result, flags)
		return
	}

	command := args[0]
	_ = ""
	if len(args) > 1 {
		_ = args[1]
	}

	// Handle help command first (no config required)
	if command == "help" {
		flags := &jira.Flags{
			Clip: toolConfig.Clip,
			File: toolConfig.File,
			JSON: toolConfig.JSON,
			Open: toolConfig.Open,
		}
		result := jira.RouteCommand("help", []string{}, nil, flags)
		jira.OutputJSON(result, flags)
		return
	}

	// Load configuration for other commands
	jiraConfig, apiToken, err := jira.LoadJiraConfig()
	jira.HandleError(err, "Failed to load configuration")

	// Create Jira client
	client := jira.NewClient(
		jiraConfig.BaseURL,
		jiraConfig.Email,
		apiToken,
		jiraConfig.DefaultProject,
	)

	// Create flags struct
	flags := &jira.Flags{
		Clip: toolConfig.Clip,
		File: toolConfig.File,
		JSON: toolConfig.JSON,
		Open: toolConfig.Open,
	}

	// Route command to appropriate handler
	remainingArgs := args[1:]
	result := jira.RouteCommand(command, remainingArgs, client, flags)
	jira.OutputJSON(result, flags)
}

func parseFlags() ToolConfig {
	config := ToolConfig{}

	flag.BoolVar(&config.Clip, "clip", false, "Copy to clipboard")
	flag.StringVar(&config.File, "file", "", "Write to file")
	flag.BoolVar(&config.JSON, "json", false, "Output in JSON format")
	flag.BoolVar(&config.Open, "o", false, "Open in browser")
	flag.BoolVar(&config.Open, "open", false, "Open in browser")

	flags.ReorderAndParse()

	return config
}

type ToolConfig struct {
	Clip bool
	File string
	JSON bool
	Open bool
}
