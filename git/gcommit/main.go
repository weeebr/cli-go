package main

// DESCRIPTION: create commit message w/ AI

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/io"
)

type Config struct {
	Compact bool
	JSON    bool
}

func main() {

	config := parseFlags()

	result, err := ai.CreateAICommit()
	ai.ExitIf(err, "failed to commit changes")

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
		if committed, exists := resultMap["committed"]; exists {
			if pushed, exists := resultMap["pushed"]; exists {
				if message, exists := resultMap["message"]; exists {
					if committed.(bool) {
						fmt.Printf("✅ Committed: %s", message)
						if pushed.(bool) {
							fmt.Printf(" and pushed\n")
						} else {
							fmt.Printf(" (not pushed)\n")
						}
					} else {
						fmt.Printf("❌ Failed to commit: %s\n", message)
					}
				}
			}
		}
	}
}
