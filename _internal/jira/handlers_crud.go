package jira

import (
	"fmt"
	"cli-go/_internal/config"
	"os/exec"
	"runtime"
	"strings"
)

// handleHistory handles history command
func handleHistory(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "history",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue key from args
	issueKey := ""
	if len(args) > 0 {
		issueKey = args[0]
	} else {
		return JiraResult{
			Action:  "history",
			Success: false,
			Error:   "Issue key required for history command",
		}
	}

	// Fetch issue with changelog
	issue, err := client.GetIssueWithChangelog(issueKey)
	if err != nil {
		return JiraResult{
			Action:  "history",
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch issue: %v", err),
		}
	}

	// Format changelog
	formatted, err := FormatChangelog(issue)
	if err != nil {
		return JiraResult{
			Action:  "history",
			Success: false,
			Error:   fmt.Sprintf("Failed to format changelog: %v", err),
		}
	}

	return JiraResult{
		Action:  "history",
		Success: true,
		Message: formatted,
	}
}

// handleCreate handles create command
func handleCreate(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "create",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue details from args
	if len(args) < 2 {
		return JiraResult{
			Action:  "create",
			Success: false,
			Error:   "Usage: jira create <project> <summary> [description]",
		}
	}

	project := args[0]
	summary := args[1]
	description := ""
	if len(args) > 2 {
		description = strings.Join(args[2:], " ")
	}

	// Create issue
	issue, err := client.CreateIssue(project, summary, description)
	if err != nil {
		return JiraResult{
			Action:  "create",
			Success: false,
			Error:   fmt.Sprintf("Failed to create issue: %v", err),
		}
	}

	// Format result
	formatted := fmt.Sprintf("✅ Created issue: %s\n%s", issue.Key, issue.Fields.Summary)
	if issue.Fields.Description != nil {
		desc, _ := ConvertADFToMarkdown(issue.Fields.Description)
		if desc != "No description" {
			formatted += "\n\n" + desc
		}
	}

	return JiraResult{
		Action:  "create",
		Success: true,
		Message: formatted,
	}
}

// handleUpdate handles update command
func handleUpdate(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "update",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue key and field from args
	if len(args) < 3 {
		return JiraResult{
			Action:  "update",
			Success: false,
			Error:   "Usage: jira update <issue-key> <field> <value>",
		}
	}

	issueKey := args[0]
	field := args[1]
	value := strings.Join(args[2:], " ")

	// Update issue
	_, err := client.UpdateIssue(issueKey, field, value)
	if err != nil {
		return JiraResult{
			Action:  "update",
			Success: false,
			Error:   fmt.Sprintf("Failed to update issue: %v", err),
		}
	}

	formatted := fmt.Sprintf("✅ Updated %s: %s", issueKey, field)
	return JiraResult{
		Action:  "update",
		Success: true,
		Message: formatted,
	}
}

// handleComment handles comment command
func handleComment(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "comment",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue key and comment from args
	if len(args) < 2 {
		return JiraResult{
			Action:  "comment",
			Success: false,
			Error:   "Usage: jira comment <issue-key> <comment>",
		}
	}

	issueKey := args[0]
	comment := strings.Join(args[1:], " ")

	// Add comment
	err := client.AddComment(issueKey, comment)
	if err != nil {
		return JiraResult{
			Action:  "comment",
			Success: false,
			Error:   fmt.Sprintf("Failed to add comment: %v", err),
		}
	}

	formatted := fmt.Sprintf("✅ Added comment to %s", issueKey)
	return JiraResult{
		Action:  "comment",
		Success: true,
		Message: formatted,
	}
}

// handleAssign handles assign command
func handleAssign(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "assign",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue key and assignee from args
	if len(args) < 2 {
		return JiraResult{
			Action:  "assign",
			Success: false,
			Error:   "Usage: jira assign <issue-key> <assignee>",
		}
	}

	issueKey := args[0]
	assignee := args[1]

	// Assign issue
	err := client.AssignIssue(issueKey, assignee)
	if err != nil {
		return JiraResult{
			Action:  "assign",
			Success: false,
			Error:   fmt.Sprintf("Failed to assign issue: %v", err),
		}
	}

	formatted := fmt.Sprintf("✅ Assigned %s to %s", issueKey, assignee)
	return JiraResult{
		Action:  "assign",
		Success: true,
		Message: formatted,
	}
}

// handleTransition handles transition command
func handleTransition(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "transition",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue key and transition from args
	if len(args) < 2 {
		return JiraResult{
			Action:  "transition",
			Success: false,
			Error:   "Usage: jira transition <issue-key> <transition>",
		}
	}

	issueKey := args[0]
	transition := args[1]

	// Transition issue
	err := client.TransitionIssue(issueKey, transition)
	if err != nil {
		return JiraResult{
			Action:  "transition",
			Success: false,
			Error:   fmt.Sprintf("Failed to transition issue: %v", err),
		}
	}

	formatted := fmt.Sprintf("✅ Transitioned %s to %s", issueKey, transition)
	return JiraResult{
		Action:  "transition",
		Success: true,
		Message: formatted,
	}
}

// handleOpen handles open command
func handleOpen(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "open",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get issue key from args
	if len(args) < 1 {
		return JiraResult{
			Action:  "open",
			Success: false,
			Error:   "Usage: jira open <issue-key>",
		}
	}

	issueKey := args[0]

	// Get base URL from config
	config, err := config.LoadConfig()
	if err != nil {
		return JiraResult{
			Action:  "open",
			Success: false,
			Error:   fmt.Sprintf("Failed to load config: %v", err),
		}
	}

	baseURL := config.Jira.BaseURL
	if baseURL == "" {
		baseURL = "https://your-company.atlassian.net"
	}

	// Construct URL
	url := fmt.Sprintf("%s/browse/%s", baseURL, issueKey)

	// Open in browser
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	err = cmd.Run()
	if err != nil {
		return JiraResult{
			Action:  "open",
			Success: false,
			Error:   fmt.Sprintf("Failed to open browser: %v", err),
		}
	}

	formatted := fmt.Sprintf("✅ Opened %s in browser", url)
	return JiraResult{
		Action:  "open",
		Success: true,
		Message: formatted,
	}
}
