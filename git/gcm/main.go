package main

// DESCRIPTION: checkout to forkpoint branch

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/git"
	"cli-go/_internal/io"
)

type Config struct {
	Compact bool
	JSON    bool
}

type CheckoutResult struct {
	Switched bool   `json:"switched"`
	Branch   string `json:"branch"`
	Stashed  bool   `json:"stashed"`
}

func main() {

	config := parseFlags()

	// Get fork-point branch
	baseBranch, err := git.GetForkPointBranch()
	ai.ExitIf(err, "failed to get fork-point branch")

	// Check if there are uncommitted changes
	modified, err := git.GetModifiedFiles()
	ai.ExitIf(err, "failed to get modified files")

	stashed := false
	if len(modified) > 0 {
		// Auto-stash changes
		_, err := git.CreateSmartStash()
		ai.ExitIf(err, "failed to create stash")
		stashed = true
	}

	// Checkout to fork-point branch
	result, err := git.CheckoutBranch(baseBranch)
	ai.ExitIf(err, "failed to checkout branch")

	response := CheckoutResult{
		Switched: result.Switched,
		Branch:   result.Branch,
		Stashed:  stashed,
	}

	if config.JSON {
		io.DirectOutput(response, *clip, *file, true)
	} else {
		outputDefault(response)
	}
}

var (
	clip = flag.Bool("clip", false, "Copy to clipboard")
	file = flag.String("file", "", "Write to file")
)

func parseFlags() Config {
	config := Config{}

	flag.BoolVar(&config.Compact, "compact", false, "Use compact JSON format")
	flag.BoolVar(&config.JSON, "json", false, "Output in JSON format")

	flag.Parse()

	return config
}

func outputDefault(response CheckoutResult) {
	if response.Switched {
		fmt.Printf("✅ Switched to branch: %s", response.Branch)
		if response.Stashed {
			fmt.Printf(" (with auto-stash)")
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("❌ Failed to switch to branch: %s\n", response.Branch)
	}
}
