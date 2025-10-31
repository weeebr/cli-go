package jira

import (
	"fmt"
	"strings"
	"time"

	"cli-go/_internal/io"
)

// FormatIssueDisplay formats an issue for display
func FormatIssueDisplay(issue *Issue, showComments, showTesting bool, baseURL string) (string, error) {
	var output strings.Builder

	// Create boxed header with auto-emoji
	title := fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary)
	boxedTitle := io.FormatBoxed(io.FormatWithEmoji(title, "ticket"))
	output.WriteString(boxedTitle + "\n\n")

	// Status line
	status := "Unknown"
	assignee := "Unassigned"
	reporter := "Unknown"

	if issue.Fields.Status.Name != "" {
		status = issue.Fields.Status.Name
	}
	if issue.Fields.Assignee.DisplayName != "" {
		assignee = issue.Fields.Assignee.DisplayName
	}
	if issue.Fields.Reporter.DisplayName != "" {
		reporter = issue.Fields.Reporter.DisplayName
	}

	output.WriteString(fmt.Sprintf("Status: %s | %s | %s\n\n",
		status,
		io.FormatWithEmoji(fmt.Sprintf("Assignee: %s", assignee), "user"),
		io.FormatWithEmoji(fmt.Sprintf("Reporter: %s", reporter), "user")))

	// Description
	if issue.Fields.Description != nil {
		description, err := ConvertADFToMarkdown(issue.Fields.Description)
		if err != nil {
			return "", fmt.Errorf("failed to convert description: %v", err)
		}
		if description != "No description" {
			output.WriteString(description + "\n")
		}
	}

	// Comments section
	if showComments && len(issue.Fields.Comments.Comments) > 0 {
		output.WriteString("## Comments\n\n")
		for i, comment := range issue.Fields.Comments.Comments {
			// Format date
			created, err := formatDate(comment.Created)
			if err != nil {
				created = comment.Created
			}

			output.WriteString(fmt.Sprintf("**%s** (%s)\n", comment.Author.DisplayName, created))

			// Convert comment body
			body, err := ConvertADFToMarkdown(comment.Body)
			if err != nil {
				body = fmt.Sprintf("Error converting comment: %v", err)
			}
			if body != "" {
				output.WriteString(body + "\n")
			}

			if i < len(issue.Fields.Comments.Comments)-1 {
				output.WriteString("\n")
			}
		}
		output.WriteString("\n")
	}

	// Testing section
	if showTesting {
		testingInstructions := formatTestingInstructions(issue)
		if testingInstructions != "No specific testing instructions found." {
			output.WriteString("## Testing Instructions\n\n")
			output.WriteString(testingInstructions)
			output.WriteString("\n")
		}
	}

	// Footer with divider and link
	output.WriteString("\n----\n")
	url := getIssueURL(baseURL, issue.Key)
	output.WriteString(io.FormatWithEmoji(url, "url"))

	return output.String(), nil
}

// FormatCommentsOnly formats only comments for display
func FormatCommentsOnly(issue *Issue) (string, error) {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("📋 Comments for %s:\n\n", issue.Key))

	if len(issue.Fields.Comments.Comments) > 0 {
		for i, comment := range issue.Fields.Comments.Comments {
			// Format date
			created, err := formatDate(comment.Created)
			if err != nil {
				created = comment.Created
			}

			output.WriteString(fmt.Sprintf("**%s** (%s)\n", comment.Author.DisplayName, created))

			// Convert comment body
			body, err := ConvertADFToMarkdown(comment.Body)
			if err != nil {
				body = fmt.Sprintf("Error converting comment: %v", err)
			}
			if body != "" {
				output.WriteString(body + "\n")
			}

			if i < len(issue.Fields.Comments.Comments)-1 {
				output.WriteString("\n")
			}
		}
	} else {
		output.WriteString("No comments.\n")
	}

	return output.String(), nil
}

// FormatTestingOnly formats only testing instructions for display
func FormatTestingOnly(issue *Issue) (string, error) {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("# %s: %s (Testing Instructions)\n\n", issue.Key, issue.Fields.Summary))

	testingInstructions := formatTestingInstructions(issue)
	output.WriteString(testingInstructions)

	return output.String(), nil
}

// FormatIssueOpenMode formats an issue for open mode (header + footer only)
func FormatIssueOpenMode(issue *Issue, baseURL string) (string, error) {
	var output strings.Builder

	// Create boxed header with auto-emoji
	title := fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary)
	boxedTitle := io.FormatBoxed(io.FormatWithEmoji(title, "ticket"))
	output.WriteString(boxedTitle + "\n\n")

	// Status line
	status := "Unknown"
	assignee := "Unassigned"
	reporter := "Unknown"

	if issue.Fields.Status.Name != "" {
		status = issue.Fields.Status.Name
	}
	if issue.Fields.Assignee.DisplayName != "" {
		assignee = issue.Fields.Assignee.DisplayName
	}
	if issue.Fields.Reporter.DisplayName != "" {
		reporter = issue.Fields.Reporter.DisplayName
	}

	output.WriteString(fmt.Sprintf("Status: %s | %s | %s\n\n",
		status,
		io.FormatWithEmoji(fmt.Sprintf("Assignee: %s", assignee), "user"),
		io.FormatWithEmoji(fmt.Sprintf("Reporter: %s", reporter), "user")))

	// Footer with divider and link
	output.WriteString("----\n")
	url := getIssueURL(baseURL, issue.Key)
	output.WriteString(io.FormatWithEmoji(url, "url"))

	return output.String(), nil
}

// AddFooter adds a footer with issue link
func AddFooter(issueKey, baseURL string) string {
	url := getIssueURL(baseURL, issueKey)
	return fmt.Sprintf("\n\n--------\n%s", io.FormatWithEmoji(url, "url"))
}

// formatDate formats an ISO date string to a readable format
func formatDate(dateStr string) (string, error) {
	// Parse ISO 8601 date
	t, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr)
	if err != nil {
		// Try without milliseconds
		t, err = time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			return dateStr, err
		}
	}

	// Format as readable date
	return t.Format("Jan 2, 2006 at 3:04 PM"), nil
}

// formatTestingInstructions extracts and formats testing instructions from custom fields
func formatTestingInstructions(issue *Issue) string {
	// Check custom fields for testing instructions
	if issue.Fields.CustomField10087 != nil {
		// Convert ADF to markdown if needed
		instructions, err := ConvertADFToMarkdown(issue.Fields.CustomField10087)
		if err == nil && instructions != "" && instructions != "No description" {
			return instructions
		}
	}

	// Check other custom fields
	if issue.Fields.CustomField10093 != nil {
		instructions, err := ConvertADFToMarkdown(issue.Fields.CustomField10093)
		if err == nil && instructions != "" && instructions != "No description" {
			return instructions
		}
	}

	if issue.Fields.CustomField10077 != nil {
		instructions, err := ConvertADFToMarkdown(issue.Fields.CustomField10077)
		if err == nil && instructions != "" && instructions != "No description" {
			return instructions
		}
	}

	return "No specific testing instructions found."
}

// getIssueURL constructs the full URL for an issue
func getIssueURL(baseURL, issueKey string) string {
	return fmt.Sprintf("%s/browse/%s", baseURL, issueKey)
}

// FormatUserActivity formats user activity report
func FormatUserActivity(username, baseURL string, viewed, created, updated *SearchResults) (string, error) {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("👤 User Activity: %s\n\n", username))

	// Viewed issues
	if viewed != nil && len(viewed.Issues) > 0 {
		output.WriteString("## Recently Viewed\n\n")
		for i, issue := range viewed.Issues {
			output.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, issue.Key, issue.Fields.Summary))
		}
		output.WriteString("\n")
	}

	// Created issues
	if created != nil && len(created.Issues) > 0 {
		output.WriteString("## Created Issues\n\n")
		for i, issue := range created.Issues {
			output.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, issue.Key, issue.Fields.Summary))
		}
		output.WriteString("\n")
	}

	// Updated issues
	if updated != nil && len(updated.Issues) > 0 {
		output.WriteString("## Updated Issues\n\n")
		for i, issue := range updated.Issues {
			output.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, issue.Key, issue.Fields.Summary))
		}
		output.WriteString("\n")
	}

	// Add footer
	output.WriteString("----\n")
	output.WriteString(io.FormatWithEmoji(fmt.Sprintf("View all activity: %s/secure/Dashboard.jspa", baseURL), "url"))

	return output.String(), nil
}

// FormatChangelog formats issue changelog for display
func FormatChangelog(issue *Issue) (string, error) {
	if issue.Changelog == nil || len(issue.Changelog.Histories) == 0 {
		return "No changelog available.", nil
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("📋 Changelog for %s:\n\n", issue.Key))

	for i, history := range issue.Changelog.Histories {
		// Format date
		created, err := formatDate(history.Created)
		if err != nil {
			created = history.Created
		}

		output.WriteString(fmt.Sprintf("**%s** (%s)\n", history.Author.DisplayName, created))

		for _, item := range history.Items {
			output.WriteString(fmt.Sprintf("- %s: %s → %s\n", item.Field, item.FromString, item.ToString))
		}

		if i < len(issue.Changelog.Histories)-1 {
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}
