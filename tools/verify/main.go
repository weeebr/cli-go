package main

// DESCRIPTION: Run comprehensive verification: go vet, staticcheck, and self-tests

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"cli-go/_internal/verify"
	"os"
	"time"
)

func main() {
	var (
		clip = flag.Bool("clip", false, "Copy to clipboard")
		file = flag.String("file", "", "Write to file")
		json = flag.Bool("json", false, "Output as JSON (for piping)")
		tool = flag.String("tool", "", "Verify specific tool only")
	)
	flags.ReorderAndParse()

	startTime := time.Now()

	// Get project root
	projectRoot, err := os.Getwd()
	ai.ExitIf(err, "failed to get current directory")

	// Find all tools
	tools := verify.FindAllTools(projectRoot)
	if *tool != "" {
		// Filter to specific tool
		filtered := []string{}
		for _, t := range tools {
			if t == *tool {
				filtered = append(filtered, t)
				break
			}
		}
		if len(filtered) == 0 {
			ai.LogError("tool '%s' not found", *tool)
			os.Exit(1)
		}
		tools = filtered
	}

	// Run verification
	results := verify.RunVerification(tools, projectRoot)
	summary := verify.CreateSummary(results, startTime)

	// Output results
	io.DirectOutput(summary, *clip, *file, *json)

	// Exit with error code if any failed
	if summary.Failed > 0 {
		os.Exit(summary.Failed)
	}
}

func formatMarkdownResult(summary verify.VerificationSummary) string {
	result := fmt.Sprintf("# Verification Summary\n\n")
	result += fmt.Sprintf("**Total Tools:** %d\n", summary.Total)
	result += fmt.Sprintf("**Passed:** %d\n", summary.Passed)
	result += fmt.Sprintf("**Failed:** %d\n", summary.Failed)
	result += fmt.Sprintf("**Duration:** %s\n\n", summary.Duration)

	if summary.Failed > 0 {
		result += "## Failed Tools\n\n"
		for _, r := range summary.Results {
			if r.Status == "failed" {
				result += fmt.Sprintf("### %s\n", r.Tool)
				result += fmt.Sprintf("- **Status:** %s\n", r.Status)
				result += fmt.Sprintf("- **Duration:** %s\n", r.Duration)
				if r.Error != "" {
					result += fmt.Sprintf("- **Error:** %s\n", r.Error)
				}
				result += "**Checks:**\n"
				for _, check := range r.Checks {
					result += fmt.Sprintf("- %s\n", check)
				}
				result += "\n"
			}
		}
	}

	return result
}
