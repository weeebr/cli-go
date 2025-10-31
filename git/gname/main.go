package main

// DESCRIPTION: check repo name

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

	// Check for help command
	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		io.LogInfo("gname - Get repository name")
		io.LogInfo("Returns the basename of the current git repository root directory")
		io.LogInfo("Output: {\"name\": \"repo-name\"} or {\"error\": \"...\"}")
		return
	}

	repoInfo, err := git.GetRepoInfo()
	ai.ExitIf(err, "failed to get git branch name")

	response := map[string]string{
		"name": repoInfo.Name,
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
	if name, exists := response["name"]; exists {
		fmt.Printf("📦 Repository name: %s\n", name)
	}
}
