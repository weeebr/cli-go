package main

// DESCRIPTION: git stash (smart)

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

	stashInfo, err := git.CreateSmartStash()
	ai.ExitIf(err, "failed to get git status")

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
		if stashed, exists := resultMap["stashed"]; exists {
			if message, exists := resultMap["message"]; exists {
				if stashed.(bool) {
					fmt.Printf("💾 Stashed: %s\n", message)
				} else {
					fmt.Printf("❌ Failed to stash: %s\n", message)
				}
			}
		}
	}
}
