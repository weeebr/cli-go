package jira

import (
	"fmt"
	"cli-go/_internal/config"
	"os/exec"
	"runtime"
	"strings"
)

// handleSearch handles search command
func handleSearch(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "search",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Build JQL query
	jql := strings.Join(args, " ")
	if jql == "" {
		jql = "ORDER BY updated DESC"
	}

	// Search issues
	results, err := client.SearchIssues(jql)
	if err != nil {
		return JiraResult{
			Action:  "search",
			Success: false,
			Error:   fmt.Sprintf("Failed to search issues: %v", err),
		}
	}

	// Format results
	formatted := fmt.Sprintf("🔍 Found %d issues:\n\n", len(results.Issues))
	for i, issue := range results.Issues {
		formatted += fmt.Sprintf("%d. %s: %s\n", i+1, issue.Key, issue.Fields.Summary)
	}

	return JiraResult{
		Action:  "search",
		Success: true,
		Message: formatted,
	}
}

// handleUserActivity handles user activity command
func handleUserActivity(args []string, client *Client, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "user_activity",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Get username from args or use current user
	username := ""
	if len(args) > 0 {
		username = args[0]
	} else {
		// Get current user from config
		config, err := config.LoadConfig()
		if err != nil {
			return JiraResult{
				Action:  "user_activity",
				Success: false,
				Error:   fmt.Sprintf("Failed to load config: %v", err),
			}
		}
		username = config.Jira.Email
	}

	if username == "" {
		return JiraResult{
			Action:  "user_activity",
			Success: false,
			Error:   "Username required for user activity command",
		}
	}

	// Get user activity
	viewedResults, err := client.GetUserViewedIssues(username)
	if err != nil {
		return JiraResult{
			Action:  "user_activity",
			Success: false,
			Error:   fmt.Sprintf("Failed to get viewed issues: %v", err),
		}
	}

	createdResults, err := client.GetUserCreatedIssues(username)
	if err != nil {
		return JiraResult{
			Action:  "user_activity",
			Success: false,
			Error:   fmt.Sprintf("Failed to get created issues: %v", err),
		}
	}

	updatedResults, err := client.GetUserUpdatedIssues(username)
	if err != nil {
		return JiraResult{
			Action:  "user_activity",
			Success: false,
			Error:   fmt.Sprintf("Failed to get updated issues: %v", err),
		}
	}

	// Format user activity
	formatted, err := FormatUserActivity(username, client.BaseURL, viewedResults, createdResults, updatedResults)
	if err != nil {
		return JiraResult{
			Action:  "user_activity",
			Success: false,
			Error:   fmt.Sprintf("Failed to format user activity: %v", err),
		}
	}

	if flags.JSON {
		return JiraResult{
			Action:  "user_activity",
			Success: true,
			Message: "User activity retrieved",
			Output:  formatted,
		}
	}

	return JiraResult{
		Action:  "user_activity",
		Success: true,
		Message: formatted,
	}
}

// handleIssue handles issue lookup
func handleIssue(client *Client, issueKey, flag string, flags *Flags) JiraResult {
	if client == nil {
		return JiraResult{
			Action:  "issue",
			Success: false,
			Error:   "Client not initialized",
		}
	}

	// Handle 'o' flag - open in browser after showing basic info
	openInBrowserAfter := flag == "o" || flags.Open
	if openInBrowserAfter {
		// Fetch basic issue info
		issue, err := client.GetIssue(issueKey)
		if err != nil {
			return JiraResult{
				Action:  "issue",
				Success: false,
				Error:   fmt.Sprintf("Failed to fetch issue: %v", err),
			}
		}

		// Format basic info (header only)
		formatted, err := FormatIssueOpenMode(issue, client.BaseURL)
		if err != nil {
			return JiraResult{
				Action:  "issue",
				Success: false,
				Error:   fmt.Sprintf("Failed to format issue: %v", err),
			}
		}

		// Open in browser silently
		normalizedKey := NormalizeIssueKey(issueKey, client.DefaultProject)
		url := fmt.Sprintf("%s/browse/%s", client.BaseURL, normalizedKey)
		if err := openInBrowser(url); err != nil {
			return JiraResult{
				Action:  "issue",
				Success: false,
				Error:   fmt.Sprintf("Failed to open issue in browser: %v", err),
			}
		}

		// Return basic info for display
		return JiraResult{
			Action:  "issue",
			Success: true,
			Message: formatted,
		}
	}

	// Fetch issue with comments
	issue, err := client.GetIssueWithComments(issueKey)
	if err != nil {
		return JiraResult{
			Action:  "issue",
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch issue: %v", err),
		}
	}

	// Determine what to show based on flag
	showComments := strings.Contains(flag, "c")
	showTesting := strings.Contains(flag, "t")

	// Format issue display
	var formatted string
	var err2 error

	if flag == "cc" {
		// Comments only
		formatted, err2 = FormatCommentsOnly(issue)
	} else if flag == "t" {
		// Testing only
		formatted, err2 = FormatTestingOnly(issue)
	} else {
		// Full display with flags (description shown by default)
		formatted, err2 = FormatIssueDisplay(issue, showComments, showTesting, client.BaseURL)
	}

	if err2 != nil {
		return JiraResult{
			Action:  "issue",
			Success: false,
			Error:   fmt.Sprintf("Failed to format issue: %v", err2),
		}
	}

	if flags.JSON {
		return JiraResult{
			Action:  "issue",
			Success: true,
			Message: "Issue retrieved",
			Output:  formatted,
		}
	}

	return JiraResult{
		Action:  "issue",
		Success: true,
		Message: formatted,
	}
}

// openInBrowser opens a URL in the default browser
func openInBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Run()
}
