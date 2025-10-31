package main

// DESCRIPTION: git stash pop

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
		io.LogInfo("gsp - Git stash pop")
		io.LogInfo("Pops the most recent stash and restores changes")
		io.LogInfo("Output: {\"message\": \"stash popped\", \"popped\": true}")
		return
	}

	stashInfo, err := git.PopStash()
	ai.ExitIf(err, "failed to get git stash pop")

	if config.JSON {
		io.DirectOutput(stashInfo, *clip, *file, true)
	} else {
		outputDefault(stashInfo)
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

func outputDefault(stashInfo interface{}) {
	if resultMap, ok := stashInfo.(map[string]interface{}); ok {
		if popped, exists := resultMap["popped"]; exists {
			if message, exists := resultMap["message"]; exists {
				if popped.(bool) {
					fmt.Printf("📤 Popped stash: %s\n", message)
				} else {
					fmt.Printf("❌ Failed to pop stash: %s\n", message)
				}
			}
		}
	}
}
