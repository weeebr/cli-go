package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// SearchPRsByQuery searches for PRs using gh CLI
func SearchPRsByQuery(owner, repo, query string) ([]PR, error) {
	cmd := exec.Command("gh", "pr", "list", "--repo", fmt.Sprintf("%s/%s", owner, repo), "--search", query, "--state", "all", "--json", "number,title,url,state,createdAt,updatedAt,body,isDraft,headRefName")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search PRs: %v", err)
	}

	var prs []PR
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to parse PR data: %v", err)
	}

	// Add owner and repo info
	for i := range prs {
		prs[i].Owner = owner
		prs[i].Repo = repo
	}

	return prs, nil
}

// GetPRDetails gets detailed information about a specific PR
func GetPRDetails(owner, repo string, number int) (*PR, error) {
	cmd := exec.Command("gh", "pr", "view", fmt.Sprintf("%d", number), "--repo", fmt.Sprintf("%s/%s", owner, repo), "--json", "number,title,url,state,createdAt,updatedAt,body,isDraft,headRefName")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get PR details: %v", err)
	}

	var pr PR
	if err := json.Unmarshal(output, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse PR data: %v", err)
	}

	pr.Owner = owner
	pr.Repo = repo

	return &pr, nil
}

// GetUserOpenPRs gets open PRs for a specific user
func GetUserOpenPRs(userEmail string, owner, repo string) ([]PR, error) {
	// Search for PRs by user email in the body or title
	query := fmt.Sprintf("author:%s", userEmail)
	return SearchPRsByQuery(owner, repo, query)
}

// GetAuthenticatedUser gets the currently authenticated GitHub user
func GetAuthenticatedUser() (string, error) {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get authenticated user: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CheckGitHubAuth checks if GitHub CLI is authenticated
func CheckGitHubAuth() error {
	cmd := exec.Command("gh", "auth", "status")
	return cmd.Run()
}

// OpenInBrowser opens a URL in the default browser
func OpenInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
