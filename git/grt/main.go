package main

// DESCRIPTION: go to repo root

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/flags"
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
		io.LogInfo("grt - Get repository root")
		io.LogInfo("Returns the absolute path to the git repository root")
		io.LogInfo("Output: {\"root\": \"/path/to/repo\"} or {\"error\": \"...\"}")
		return
	}

	repoInfo, err := git.GetRepoInfo()
	ai.ExitIf(err, "failed to get git root")

	response := map[string]string{
		"root": repoInfo.Root,
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

	flags.ReorderAndParse()

	return config
}

func outputDefault(response map[string]string) {
	if root, exists := response["root"]; exists {
		fmt.Printf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Repository root: %s", root), "file"))
	}
}
