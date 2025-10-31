package main

// DESCRIPTION: displays what you see atm

import (
	"bufio"
	"cli-go/_internal/flags"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ToolInfo struct {
	Name        string
	Description string
}

type Category struct {
	Name  string
	Tools []ToolInfo
}

func main() {
	// Check for help command before parsing flags
	if len(os.Args) > 1 && os.Args[1] == "help" {
		fmt.Fprintf(os.Stderr, "help: Display formatted help for all CLI tools\n")
		fmt.Fprintf(os.Stderr, "Usage: help\n")
		fmt.Fprintf(os.Stderr, "Shows categorized list of all available CLI tools in formatted output\n")
		return
	}

	flags.ReorderAndParse()

	// Define category mapping
	categoryMap := map[string]string{
		"ai":      "AI",
		"git":     "Git",
		"tools":   "Tools",
		"ringier": "Ringier",
	}

	// Define which tools belong to which categories
	toolCategories := map[string]string{
		// AI tools
		"cld": "ai", "gem": "ai", "gro": "ai", "grop": "ai", "haik": "ai",
		"j": "ai", "ji": "ai", "jj": "ai", "jp": "ai", "prompts": "ai",

		// Git tools
		"gaff": "git", "gbd": "git", "gcb": "git", "gcd": "git", "gcm": "git",
		"gco": "git", "gcommit": "git", "ginstall": "git", "gmain": "git",
		"gname": "git", "greinstall": "git", "grt": "git", "gs": "git",
		"gsp": "git", "gstats": "git", "gprs": "git",

		// Core tools
		"check_alias": "core", "killport": "core", "perf": "core",
		"zsh": "core", "zss": "core",

		// Tools category
		"edit": "tools", "figma": "tools",
		"help": "tools", "jira": "tools",
		"repos": "tools", "test": "tools", "web": "tools",

		// Ringier
		"smart_start": "ringier",
	}

	// Add special cases for shell aliases (none currently)
	specialCases := map[string]ToolInfo{}

	// Get the directory of the current executable
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get executable path: %v\n", err)
		os.Exit(1)
	}
	
	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)
	
	// Get the cli-go directory (parent of bin/)
	goCliDir := filepath.Join(execDir, "..")

	// Scan for tools and extract descriptions
	tools := make(map[string]ToolInfo)

	// Scan ai/ directory
	scanDirectory(filepath.Join(goCliDir, "ai"), tools)

	// Scan git/ directory
	scanDirectory(filepath.Join(goCliDir, "git"), tools)

	// Scan tools/ directory
	scanDirectory(filepath.Join(goCliDir, "tools"), tools)

	// Add special cases
	for name, toolInfo := range specialCases {
		tools[name] = toolInfo
	}

	// Group tools by category
	categories := make(map[string][]ToolInfo)
	for toolName, toolInfo := range tools {
		if category, exists := toolCategories[toolName]; exists {
			categoryName := categoryMap[category]
			if categoryName == "" {
				categoryName = strings.Title(category)
			}
			categories[categoryName] = append(categories[categoryName], toolInfo)
		}
	}

	// Display formatted help with exact zsh function format
	fmt.Println()

	// Define category order
	categoryOrder := []string{"AI", "Git", "Ringier", "Core", "Tools"}

	for _, categoryName := range categoryOrder {
		if tools, exists := categories[categoryName]; exists && len(tools) > 0 {
			// Bold light blue category heading
			fmt.Printf("\033[1m%s\033[22m\n\n", categoryName)

			// Sort tools by name
			sort.Slice(tools, func(i, j int) bool {
				return tools[i].Name < tools[j].Name
			})

			total := len(tools)
			perCol := (total + 1) / 2

			// Display in 2-column format matching zsh function
			for row := 0; row < perCol; row++ {
				left := ToolInfo{}
				right := ToolInfo{}

				if row < len(tools) {
					left = tools[row]
				}
				if row+perCol < len(tools) {
					right = tools[row+perCol]
				}

				var leftCell, rightCell string

				if left.Name != "" {
					leftCell = fmt.Sprintf("\033[36m%-16s\033[0m  \033[90m%-42s\033[0m", left.Name, left.Description)
				}

				if right.Name != "" {
					rightCell = fmt.Sprintf("\033[36m%-16s\033[0m  \033[90m%s\033[0m", right.Name, right.Description)
				}

				fmt.Printf("  %-62s  %s\n", leftCell, rightCell)
			}

			fmt.Println()
		}
	}
}

func scanDirectory(dir string, tools map[string]ToolInfo) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		toolName := entry.Name()
		mainGoPath := filepath.Join(dir, toolName, "main.go")

		// Read main.go file and extract description
		description := extractDescription(mainGoPath)
		if description != "" {
			tools[toolName] = ToolInfo{
				Name:        toolName,
				Description: description,
			}
		}
	}
}

func extractDescription(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "// DESCRIPTION:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "// DESCRIPTION:"))
		}
	}

	return ""
}
