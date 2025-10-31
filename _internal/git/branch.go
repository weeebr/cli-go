package git

import (
	"fmt"
	"cli-go/_internal/sys"
	"strings"
	"time"
)

// CheckoutResult holds the result of a checkout operation
type CheckoutResult struct {
	Switched bool   `json:"switched"`
	Branch   string `json:"branch"`
	Error    string `json:"error,omitempty"`
}

// BranchOperationResult represents the result of a branch operation
type BranchOperationResult struct {
	Success bool   `json:"success"`
	Branch  string `json:"branch"`
	Error   string `json:"error,omitempty"`
}

// Branch represents a Git branch with metadata
type Branch struct {
	Name       string    `json:"name"`
	LastCommit time.Time `json:"lastCommit"`
	Author     string    `json:"author"`
	RepoPath   string    `json:"repoPath"`
	HasPR      bool      `json:"hasPR"`
	PRNumber   int       `json:"prNumber,omitempty"`
}

// GetMainBranch returns the main branch name (main, master, etc.)
func GetMainBranch() (string, error) {
	// Check common main branch names in order of preference
	branches := []string{"main", "trunk", "mainline", "default", "stable", "master"}

	for _, branch := range branches {
		result := sys.RunCommand("git", "show-ref", "-q", "--verify", "refs/heads/"+branch)
		if result.ExitCode == 0 {
			return branch, nil
		}
	}

	// Check remote branches
	for _, branch := range branches {
		result := sys.RunCommand("git", "show-ref", "-q", "--verify", "refs/remotes/origin/"+branch)
		if result.ExitCode == 0 {
			return branch, nil
		}
	}

	// Check remote HEAD
	result := sys.RunCommand("git", "rev-parse", "--abbrev-ref", "origin/HEAD")
	if result.ExitCode == 0 && strings.HasPrefix(result.Stdout, "origin/") {
		return strings.TrimPrefix(result.Stdout, "origin/"), nil
	}

	// Default to master
	return "master", nil
}

// GetBranchCount returns the number of local branches
func GetBranchCount(repoPath string) int {
	result := sys.RunCommandInDir(repoPath, "git", "branch")
	if result.ExitCode != 0 {
		return 0
	}
	return len(strings.Split(result.Stdout, "\n")) - 1
}

// ListBranches returns all branches in the repository
func ListBranches(repoPath string) ([]string, error) {
	result := sys.RunCommand("git", "-C", repoPath, "branch", "-a")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get branches: %s", result.Stderr)
	}

	var branches []string
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove * prefix and remote/ prefix
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimPrefix(line, "remotes/origin/")
		line = strings.TrimPrefix(line, "origin/")

		if line != "" && !strings.Contains(line, "HEAD") {
			branches = append(branches, line)
		}
	}

	return branches, nil
}

// GetUserBranches returns branches created by a specific author
func GetUserBranches(repoPath, authorEmail string) ([]Branch, error) {
	// Get all branches
	branches, err := ListBranches(repoPath)
	if err != nil {
		return nil, err
	}

	var userBranches []Branch
	for _, branchName := range branches {
		// Get last commit info for this branch
		lastCommit, author, err := GetBranchLastCommit(repoPath, branchName)
		if err != nil {
			continue // Skip branches we can't analyze
		}

		// Check if this branch was created by the user
		if strings.Contains(strings.ToLower(author), strings.ToLower(authorEmail)) {
			userBranches = append(userBranches, Branch{
				Name:       branchName,
				LastCommit: lastCommit,
				Author:     author,
				RepoPath:   repoPath,
				HasPR:      false, // Will be determined separately
			})
		}
	}

	return userBranches, nil
}

// GetBranchLastCommit returns the last commit info for a branch
func GetBranchLastCommit(repoPath, branch string) (time.Time, string, error) {
	// Get last commit date and author
	result := sys.RunCommand("git", "-C", repoPath, "log", "-1", "--format=%ad|%ae", "--date=iso", branch)
	if result.ExitCode != 0 {
		return time.Time{}, "", fmt.Errorf("failed to get last commit for branch %s: %s", branch, result.Stderr)
	}

	parts := strings.Split(strings.TrimSpace(result.Stdout), "|")
	if len(parts) != 2 {
		return time.Time{}, "", fmt.Errorf("unexpected commit format")
	}

	dateStr := parts[0]
	author := parts[1]

	// Parse the date
	date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
	if err != nil {
		// Try alternative format
		date, err = time.Parse("2006-01-02T15:04:05-07:00", dateStr)
		if err != nil {
			return time.Time{}, "", fmt.Errorf("failed to parse date: %v", err)
		}
	}

	return date, author, nil
}

// CheckBranchHasPR checks if a branch has an associated PR (placeholder for GitHub API integration)
func CheckBranchHasPR(repoPath, branch string) (bool, int, error) {
	// This would need GitHub API integration to check for PRs
	// For now, return false - this will be implemented when integrating with gh API
	return false, 0, nil
}

// DeleteBranch deletes a branch
func DeleteBranch(branch string) (*BranchOperationResult, error) {
	result := sys.RunCommand("git", "branch", "--delete", branch)
	if result.ExitCode != 0 {
		return &BranchOperationResult{
			Success: false,
			Branch:  branch,
			Error:   result.Stderr,
		}, fmt.Errorf("failed to delete branch %s: %s", branch, result.Stderr)
	}

	return &BranchOperationResult{
		Success: true,
		Branch:  branch,
	}, nil
}
