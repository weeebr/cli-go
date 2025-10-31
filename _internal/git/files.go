package git

import (
	"fmt"
	"cli-go/_internal/sys"
	"strings"
)

// GetStagedFiles returns a list of staged files
func GetStagedFiles() ([]string, error) {
	result := sys.RunCommand("git", "diff", "--staged", "--name-only")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get staged files: %s", result.Stderr)
	}

	if result.Stdout == "" {
		return []string{}, nil
	}

	return strings.Split(result.Stdout, "\n"), nil
}

// GetModifiedFiles returns a list of modified files
func GetModifiedFiles() ([]string, error) {
	result := sys.RunCommand("git", "diff", "--name-only")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get modified files: %s", result.Stderr)
	}

	if result.Stdout == "" {
		return []string{}, nil
	}

	return strings.Split(result.Stdout, "\n"), nil
}

// GetUntrackedFiles returns a list of untracked files
func GetUntrackedFiles() ([]string, error) {
	result := sys.RunCommand("git", "ls-files", "--others", "--exclude-standard")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get untracked files: %s", result.Stderr)
	}

	if result.Stdout == "" {
		return []string{}, nil
	}

	return strings.Split(result.Stdout, "\n"), nil
}

// GetTrackedFiles returns a list of all git-tracked files
func GetTrackedFiles(repoPath string) ([]string, error) {
	result := sys.RunCommandInDir(repoPath, "git", "ls-files")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get tracked files: %s", result.Stderr)
	}

	if result.Stdout == "" {
		return []string{}, nil
	}

	return strings.Split(result.Stdout, "\n"), nil
}
