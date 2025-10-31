package ai

import (
	"fmt"
	"cli-go/_internal/sys"
)

// CommitResult holds the result of an AI commit
type CommitResult struct {
	Committed bool   `json:"committed"`
	Pushed    bool   `json:"pushed"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
}

// GenerateCommitMessage generates a commit message using AI
func GenerateCommitMessage() (string, error) {
	// Get staged changes
	result := sys.RunCommand("git", "diff", "--staged")
	if result.ExitCode != 0 {
		return "", fmt.Errorf("failed to get staged changes: %s", result.Stderr)
	}

	if result.Stdout == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	// Use ChatGPT to generate commit message
	prompt := `generate a commit message with these exact requirements:
- maximum 70 characters
- all lowercase letters only
- must start with add, remove, update, or improve
- no backticks, colons, or special characters
- either one concise sentence or two parts joined with +
- base only on the actual changes shown in the git diff
- do not make assumptions beyond what is explicitly shown`

	// Call ChatGPT (assuming it's available in PATH)
	chatResult := sys.RunCommand("chatgpt", prompt)
	if chatResult.ExitCode != 0 {
		return "", fmt.Errorf("failed to generate commit message: %s", chatResult.Stderr)
	}

	return chatResult.Stdout, nil
}

// CreateAICommit creates a commit with AI-generated message
func CreateAICommit() (*CommitResult, error) {
	// Check if there are staged changes
	stagedResult := sys.RunCommand("git", "diff", "--staged", "--name-only")
	if stagedResult.ExitCode != 0 {
		return &CommitResult{
			Committed: false,
			Error:     "not a git repository",
		}, fmt.Errorf("not a git repository")
	}

	if stagedResult.Stdout == "" {
		// No staged changes, check for unstaged changes
		unstagedResult := sys.RunCommand("git", "diff", "--name-only")
		if unstagedResult.ExitCode != 0 || unstagedResult.Stdout == "" {
			return &CommitResult{
				Committed: false,
				Error:     "no changes to commit",
			}, fmt.Errorf("no changes to commit")
		}

		// Stage all changes
		stageResult := sys.RunCommand("git", "add", "-A")
		if stageResult.ExitCode != 0 {
			return &CommitResult{
				Committed: false,
				Error:     "failed to stage changes",
			}, fmt.Errorf("failed to stage changes: %s", stageResult.Stderr)
		}
	}

	// Generate commit message
	message, err := GenerateCommitMessage()
	if err != nil {
		return &CommitResult{
			Committed: false,
			Error:     err.Error(),
		}, err
	}

	// Create commit
	commitResult := sys.RunCommand("git", "commit", "-m", message)
	if commitResult.ExitCode != 0 {
		return &CommitResult{
			Committed: false,
			Error:     "failed to commit",
		}, fmt.Errorf("failed to commit: %s", commitResult.Stderr)
	}

	// Push changes
	pushResult := sys.RunCommand("git", "push")
	pushed := pushResult.ExitCode == 0

	return &CommitResult{
		Committed: true,
		Pushed:    pushed,
		Message:   message,
	}, nil
}
