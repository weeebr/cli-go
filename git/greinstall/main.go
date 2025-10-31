package main

// DESCRIPTION: reinstall repo

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"cli-go/_internal/sys"
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
		io.LogInfo("greinstall - Reinstall repository packages")
		io.LogInfo("Removes node_modules and reinstalls packages")
		io.LogInfo("Output: {\"reinstalled\": true, \"manager\": \"pnpm\"}")
		return
	}

	result := sys.ReinstallPackages()
	if !result.Installed {
		ai.ExitIf(fmt.Errorf("reinstallation failed: %s", result.Error), "package reinstallation failed")
	}

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

	flags.ReorderAndParse()

	return config
}

func outputDefault(result interface{}) {
	if resultMap, ok := result.(map[string]interface{}); ok {
		if installed, exists := resultMap["installed"]; exists {
			if manager, exists := resultMap["manager"]; exists {
				if installed.(bool) {
					fmt.Printf("🔄 Reinstalled packages using %s\n", manager)
				} else {
					fmt.Printf("❌ Failed to reinstall packages using %s\n", manager)
				}
			}
		}
	}
}
