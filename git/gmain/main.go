package main

// DESCRIPTION: check main branch name

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

func main() {

	config := parseFlags()

	branch, err := git.GetMainBranch()
	ai.ExitIf(err, "failed to get main branch")

	response := map[string]string{
		"branch": branch,
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

func outputDefault(response map[string]string) {
	if branch, exists := response["branch"]; exists {
		fmt.Printf("🌿 Main branch: %s\n", branch)
	}
}
