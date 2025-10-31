package jira

import (
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/io"
	"os"
)

// Flags represents the command line flags for jira tool
type Flags struct {
	Clip bool
	File string
	JSON bool
	Open bool
}

// JiraResult represents the result of a jira command
type JiraResult struct {
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
	Output  string `json:"output,omitempty"`
}

// RouteCommand routes commands to appropriate handlers
func RouteCommand(command string, args []string, client *Client, flags *Flags) JiraResult {
	switch command {
	case "history":
		return handleHistory(args, client, flags)
	case "u", "user":
		return handleUserActivity(args, client, flags)
	case "help":
		return handleHelp(flags)
	default:
		// Default: treat as issue key
		// Extract flag from args if present
		flag := ""
		if len(args) > 0 {
			flag = args[0]
		}
		return handleIssue(client, command, flag, flags)
	}
}

// handleHelp handles help command
func handleHelp(flags *Flags) JiraResult {
	helpText := `Jira CLI Tool

Usage:
  jira <issue-key> [flags]     - View issue details
  jira history <issue-key>     - View issue history
  jira user [username]         - View user activity
  jira help                    - Show this help

Flags:
  -c    Show comments
  -t    Show testing instructions
  -o    Open in browser (minimal display)
  -cc   Comments only
  -json Output as JSON
  -compact Compact output format

Examples:
  jira PNT-123
  jira PNT-123 c
  jira PNT-123 t
  jira PNT-123 o
  jira history PNT-123
  jira user
  jira user john.doe@company.com
`

	return JiraResult{
		Action:  "help",
		Success: true,
		Message: helpText,
	}
}

// HandleError handles errors and exits if needed
func HandleError(err error, message string) {
	ai.ExitIf(err, message)
}

// OutputJSON outputs result as JSON if requested
func OutputJSON(result JiraResult, flags *Flags) {
	if flags.JSON {
		// Use DirectOutput for consistent behavior
		io.DirectOutput(result, flags.Clip, flags.File, flags.JSON)
	} else {
		if result.Success {
			fmt.Print(result.Message)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
		}
	}
}
