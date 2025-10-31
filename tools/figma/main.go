package main

// DESCRIPTION: Figma component search and management tool

import (
	"flag"
	"fmt"
	"os"

	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/figma"
)

const defaultFileKey = "Bvw817OVY6zhmEty1Syj8Q" // ORBIT_FILE_KEY

func main() {
	var (
		clip    = flag.Bool("clip", false, "Copy to clipboard")
		file    = flag.String("file", "", "Write to file")
		compact = flag.Bool("compact", false, "Compact JSON output")
		json    = flag.Bool("json", false, "Output as JSON (for piping)")
	)

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: command or query required\n")
		fmt.Fprintf(os.Stderr, "Usage: figma <query> | search | init | list | cache clear | --\n")
		os.Exit(1)
	}

	// Create store once for commands that need credentials

	// Route to command handlers
	command := args[0]
	switch command {
	case "search":
		token := getFigmaToken()
		figma.HandleSearch(args[1:], token, *clip, *file, *compact, *json)
	case "--":
		token := getFigmaToken()
		figma.HandleFullMetadata(args[1:], token, *clip, *file, *compact, *json)
	case "init":
		token := getFigmaToken()
		figma.HandleInit(args[1:], token, *clip, *file, *compact, *json)
	case "list":
		token := getFigmaToken()
		figma.HandleList(args[1:], token, *clip, *file, *compact, *json)
	case "cache":
		token := getFigmaToken()
		figma.HandleCache(args[1:], token, *clip, *file, *compact, *json)
	default:
		// Default: treat as search query
		token := getFigmaToken()
		figma.HandleSearch(args, token, *clip, *file, *compact, *json)
	}
}

func getFigmaToken() string {
	token, err := config.GetKey("figma")
	ai.ExitIf(err, "failed to get Figma API token")
	return token
}

func showHelp() {
	help := `Figma Component Search and Management Tool

USAGE:
    figma [OPTIONS] [COMMAND] [ARGS...]

COMMANDS:
    help                    Show this help message
    search <query>          Search for components by name
    init [fileKey]          Initialize cache with components from file
    list                    List all cached components
    cache <stats|clear>     Manage cache
    -- <query>              Get full metadata for component

OPTIONS:
    --clip                                    Copy to clipboard
    --file=./output.json                      Write to file
    --compact                                 Use compact JSON format
    --json                                    Output in JSON format

EXAMPLES:
    figma search "button"                     Search for button components
    figma init                                Initialize cache with default file
    figma list                                List all cached components
    figma cache stats                         Show cache statistics
    figma cache clear                         Clear cache
    figma -- "button-main"                 Get full metadata for component

For more information, visit: https://github.com/your-repo/figma-cli`

	fmt.Println(help)
}
