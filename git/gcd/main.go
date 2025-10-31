package main

// DESCRIPTION: git checkout develop

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

	result, err := git.CheckoutBranch("develop")
	ai.ExitIf(err, "failed to change directory")

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
		if switched, exists := resultMap["switched"]; exists {
			if branch, exists := resultMap["branch"]; exists {
				if switched.(bool) {
					fmt.Printf("✅ Switched to branch: %s\n", branch)
				} else {
					fmt.Printf("❌ Failed to switch to branch: %s\n", branch)
				}
			}
		}
	}
}
