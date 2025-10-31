package main

// DESCRIPTION: git branch --delete $1

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/git"
	"cli-go/_internal/io"
	"os"
)

type Config struct {
	Compact bool
	JSON    bool
}

type Input struct {
	Branch string `json:"branch"`
}

func main() {

	config := parseFlags()

	// Read input from stdin
	var input Input
	if err := io.ReadJSON(&input); err != nil {
		ai.ExitIf(err, "failed to read input")
	}

	if input.Branch == "" {
		ai.LogError("branch is required")
		os.Exit(1)
	}

	result, err := git.DeleteBranch(input.Branch)
	ai.ExitIf(err, "failed to delete branch")

	if config.JSON {
		io.DirectOutput(result, *clip, *file, true)
	} else {
		outputDefault(result)
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

func outputDefault(result interface{}) {
	if resultMap, ok := result.(map[string]interface{}); ok {
		if deleted, exists := resultMap["deleted"]; exists {
			if branch, exists := resultMap["branch"]; exists {
				if deleted.(bool) {
					fmt.Printf("🗑️  Deleted branch: %s\n", branch)
				} else {
					fmt.Printf("❌ Failed to delete branch: %s\n", branch)
				}
			}
		}
	}
}
