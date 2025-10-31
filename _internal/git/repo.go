package git

import (
	"fmt"
	"cli-go/_internal/sys"
	"os"
	"strings"
)

// RepoInfo holds basic repository information
type RepoInfo struct {
	Name   string `json:"name"`
	Root   string `json:"root"`
	Branch string `json:"branch"`
}

// IsGitRepo checks if the current directory is a git repository
func IsGitRepo() bool {
	result := sys.RunCommand("git", "rev-parse", "--git-dir")
	return result.ExitCode == 0
}

// IsGitRepoAtPath checks if a specific path is a git repository
func IsGitRepoAtPath(path string) bool {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// Run git command in that directory
	result := sys.RunCommandInDir(path, "git", "rev-parse", "--git-dir")
	return result.ExitCode == 0
}

// GetGitRoot returns the root directory of the git repository
func GetGitRoot() (string, error) {
	result := sys.RunCommand("git", "rev-parse", "--show-toplevel")
	if result.ExitCode != 0 {
		return "", fmt.Errorf("not a git repository: %s", result.Stderr)
	}
	return strings.TrimSpace(result.Stdout), nil
}

// GetCurrentBranch returns the current git branch
func GetCurrentBranch() (string, error) {
	result := sys.RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if result.ExitCode != 0 {
		return "", fmt.Errorf("failed to get current branch: %s", result.Stderr)
	}
	return result.Stdout, nil
}

// GetRepoName returns the repository name (basename of git root)
func GetRepoName() (string, error) {
	root, err := GetGitRoot()
	if err != nil {
		return "", err
	}

	// Get basename of the root path
	parts := strings.Split(root, "/")
	return parts[len(parts)-1], nil
}

// GetRepoInfo returns basic information about the current repository
func GetRepoInfo() (*RepoInfo, error) {
	if !IsGitRepo() {
		return nil, fmt.Errorf("not a git repository")
	}

	root, err := GetGitRoot()
	if err != nil {
		return nil, err
	}

	name, err := GetRepoName()
	if err != nil {
		return nil, err
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	return &RepoInfo{
		Name:   name,
		Root:   root,
		Branch: branch,
	}, nil
}
