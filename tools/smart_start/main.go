package main

// DESCRIPTION: Smart Start: ginstall && _smart_start

import (
	"flag"
	"fmt"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"cli-go/_internal/smart_start"
	"os"
	"path/filepath"
)


func main() {

	var (
		clip = flag.Bool("clip", false, "Copy to clipboard")
		file = flag.String("file", "", "Write to file")
		json = flag.Bool("json", false, "Output as JSON (for piping)")
	)
	flags.ReorderAndParse()

	// Get repo root from current working directory (pure Go)
	repoRoot, err := findGitRoot()
	
	// Debug: print current working directory
	if cwd, err := os.Getwd(); err == nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Current working directory: %s\n", cwd)
	}
	if err != nil {
		result := smart_start.SmartStartInternalResult{
			Action:  "_smart_start",
			Success: false,
			Error:   "not a git repository",
			Message: "Must be run from within a git repository",
		}
		if *json {
			// Use centralized output routing
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		os.Exit(1)
	}

	repoName := filepath.Base(repoRoot)

	// Check for direct project argument
	args := flag.Args()
	var directProject string
	if len(args) > 0 {
		directProject = args[0]
	}

	// Detect project type and get configuration using internal package
	result := smart_start.DetectAndConfigureProject(repoName, repoRoot, directProject)
	if *json {
		// Use centralized output routing
		io.DirectOutput(result, *clip, *file, *json)
	} else {
		outputDefault(result)
	}
}


func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}
	
	return "", fmt.Errorf("not a git repository")
}

func outputDefault(result smart_start.SmartStartInternalResult) {
	if result.Success {
		fmt.Printf("🚀 Smart Start: %s\n", result.Message)
		fmt.Printf("   Repository: %s\n", result.RepoName)
		fmt.Printf("   Project Type: %s\n", result.ProjectType)
		if result.Project != "" {
			fmt.Printf("   Project: %s\n", result.Project)
		}
		if result.Environment != "" {
			fmt.Printf("   Environment: %s\n", result.Environment)
		}
		if len(result.Config) > 0 {
			fmt.Printf("   Configuration: %v\n", result.Config)
		}
	} else {
		fmt.Printf("❌ %s\n", result.Message)
		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}
	}
}
