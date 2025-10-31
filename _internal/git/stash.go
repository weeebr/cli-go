package git

import (
	"fmt"
	"cli-go/_internal/sys"
	"strings"
	"time"
)

// StashInfo holds information about a stash
type StashInfo struct {
	Message string `json:"message"`
	Stashed bool   `json:"stashed"`
}

// CreateSmartStash creates a stash with a smart message including timestamp and file names
func CreateSmartStash() (*StashInfo, error) {
	// Get modified files (up to 3)
	modified, err := GetModifiedFiles()
	if err != nil {
		return nil, err
	}

	// Get untracked files
	untracked, err := GetUntrackedFiles()
	if err != nil {
		return nil, err
	}

	// Combine and limit to 3 files
	allFiles := append(modified, untracked...)
	if len(allFiles) > 3 {
		allFiles = allFiles[:3]
	}

	// Create message with timestamp and file names
	timestamp := time.Now().Format("15:04")
	var message string

	if len(allFiles) > 0 {
		// Get basenames of files
		basenames := make([]string, len(allFiles))
		for i, file := range allFiles {
			parts := strings.Split(file, "/")
			basenames[i] = parts[len(parts)-1]
		}
		filesStr := strings.Join(basenames, " | ")
		message = fmt.Sprintf("[%s] %s", timestamp, filesStr)
	} else {
		message = fmt.Sprintf("[%s]", timestamp)
	}

	// Create the stash
	result := sys.RunCommand("git", "stash", "push", "-m", message)
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to create stash: %s", result.Stderr)
	}

	return &StashInfo{
		Message: message,
		Stashed: true,
	}, nil
}

// PopStash pops the most recent stash
func PopStash() (*StashInfo, error) {
	result := sys.RunCommand("git", "stash", "pop")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to pop stash: %s", result.Stderr)
	}

	// Parse the stash message from the output
	message := "stash popped"
	if strings.Contains(result.Stdout, "On branch") {
		// Extract stash message if available
		lines := strings.Split(result.Stdout, "\n")
		for _, line := range lines {
			if strings.Contains(line, "WIP on") {
				message = strings.TrimSpace(line)
				break
			}
		}
	}

	return &StashInfo{
		Message: message,
		Stashed: true,
	}, nil
}

// ListStashes returns a list of available stashes
func ListStashes() ([]string, error) {
	result := sys.RunCommand("git", "stash", "list")
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to list stashes: %s", result.Stderr)
	}

	if result.Stdout == "" {
		return []string{}, nil
	}

	return strings.Split(result.Stdout, "\n"), nil
}
