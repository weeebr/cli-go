package git

import (
	"fmt"
	"cli-go/_internal/sys"
)

// CheckoutBranch checks out to a specific branch
func CheckoutBranch(branch string) (*CheckoutResult, error) {
	result := sys.RunCommand("git", "checkout", branch)
	if result.ExitCode != 0 {
		return &CheckoutResult{
			Switched: false,
			Branch:   branch,
			Error:    result.Stderr,
		}, fmt.Errorf("failed to checkout branch %s: %s", branch, result.Stderr)
	}

	return &CheckoutResult{
		Switched: true,
		Branch:   branch,
	}, nil
}

// CreateBranch creates a new branch and checks it out
func CreateBranch(branch string) (*CheckoutResult, error) {
	result := sys.RunCommand("git", "checkout", "-b", branch)
	if result.ExitCode != 0 {
		return &CheckoutResult{
			Switched: false,
			Branch:   branch,
			Error:    result.Stderr,
		}, fmt.Errorf("failed to create branch %s: %s", branch, result.Stderr)
	}

	return &CheckoutResult{
		Switched: true,
		Branch:   branch,
	}, nil
}
