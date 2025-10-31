package main

// DESCRIPTION: search for PRs by ticket ID

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/custom"
	"cli-go/_internal/flags"
	"cli-go/_internal/github"
	"cli-go/_internal/io"
	"cli-go/_internal/sys"
	"os"
	"strings"
)

type Config struct {
	Compact bool
	Single  string
	Main    bool
	All     bool
	JSON    bool
	Open    bool
}

func main() {
	cfg := parseFlags()

	// Check for help command
	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		io.LogInfo("gprs - Search for PRs by ticket ID")
		io.LogInfo("Searches for pull requests containing the specified ticket ID")
		io.LogInfo("Supports --single <path> (specific repo), --main (main repos), --all (all repos)")
		io.LogInfo("Default: all repositories")
		io.LogInfo("Flags: --json (JSON output), --o (open in browser)")
		io.LogInfo("Usage: gprs PNT-123 [flags]")
		return
	}
	if len(args) == 0 {
		io.LogError("Ticket ID required")
		os.Exit(1)
	}

	ticketID := args[0]

	// Check for 'o' flag from config
	openFlag := cfg.Open

	// Load config to get default project key
	configData, err := config.LoadConfig()
	ai.ExitIf(err, "failed to load configuration")

	// Validate ticket ID
	validatedTicketID, err := custom.ValidateTicketID(ticketID)
	ai.ExitIf(err, "invalid ticket ID")

	// Normalize ticket ID
	normalizedTicketID := custom.NormalizeTicketID(validatedTicketID, configData.Ringier.DefaultProjectKey)

	// Search for PRs using GitHub API
	var allPRs []github.PR

	// First, try main repositories only
	mainRepos := []config.RepoConfig{}
	for _, repoConfig := range configData.Repositories {
		if repoConfig.Main {
			mainRepos = append(mainRepos, repoConfig)
		}
	}

	// Search main repos first
	for _, repoConfig := range mainRepos {
		query := fmt.Sprintf("%s in:title,body", normalizedTicketID)
		prs, err := github.SearchPRsByQuery(repoConfig.Owner, repoConfig.Repo, query)
		if err != nil {
			io.LogWarning("Failed to search %s/%s: %v", repoConfig.Owner, repoConfig.Repo, err)
			continue
		}
		allPRs = append(allPRs, prs...)
	}

	// If no PRs found in main repos, search all repos
	if len(allPRs) == 0 {
		for _, repoConfig := range configData.Repositories {
			query := fmt.Sprintf("%s in:title,body", normalizedTicketID)
			prs, err := github.SearchPRsByQuery(repoConfig.Owner, repoConfig.Repo, query)
			if err != nil {
				io.LogWarning("Failed to search %s/%s: %v", repoConfig.Owner, repoConfig.Repo, err)
				continue
			}
			allPRs = append(allPRs, prs...)
		}
	}

	if len(allPRs) == 0 {
		io.LogError("No PRs found for ticket %s", normalizedTicketID)
		os.Exit(1)
	}

	// Format and output results
	output := formatPRs(allPRs, cfg.JSON)
	io.DirectOutput(output, *clip, *file, cfg.JSON)

	// If open flag is set, open the best/most relevant PR in browser (silently)
	if openFlag {
		bestPR := findBestPR(allPRs)
		if err := sys.OpenInBrowser(bestPR.URL); err != nil {
			io.LogError("Failed to open browser: %v", err)
			os.Exit(1)
		}
	}
}

var (
	clip = flag.Bool("clip", false, "Copy to clipboard")
	file = flag.String("file", "", "Write to file")
)

func parseFlags() Config {
	config := Config{}

	flag.BoolVar(&config.Compact, "compact", false, "Use compact JSON format")
	flag.StringVar(&config.Single, "single", "", "Operate on specific repository path")
	flag.BoolVar(&config.Main, "main", false, "Operate on main repositories (orbit + rasch-stack)")
	flag.BoolVar(&config.All, "all", false, "Operate on all repositories from config")
	flag.BoolVar(&config.JSON, "json", false, "Output in JSON format")
	flag.BoolVar(&config.Open, "o", false, "Open PR in browser")
	flag.BoolVar(&config.Open, "open", false, "Open PR in browser")

	flags.ReorderAndParse()

	// No default needed - GetReposToProcess handles defaults

	return config
}

// formatPRs formats a list of PRs for display
func formatPRs(prs []github.PR, jsonFormat bool) string {
	if jsonFormat {
		return formatPRsJSON(prs)
	}
	return formatPRsMarkdown(prs)
}

// formatPRsMarkdown formats PRs in markdown format
func formatPRsMarkdown(prs []github.PR) string {
	if len(prs) == 0 {
		return "No PRs found."
	}

	var result []string

	if len(prs) == 1 {
		pr := prs[0]
		result = append(result, fmt.Sprintf("## Found PR: #%d", pr.Number))
		result = append(result, fmt.Sprintf("**Title:** %s", pr.Title))
		result = append(result, fmt.Sprintf("**URL:** %s", pr.URL))
		result = append(result, fmt.Sprintf("**State:** %s", pr.State))
		result = append(result, fmt.Sprintf("**Branch:** %s", pr.HeadRefName))
		result = append(result, fmt.Sprintf("**Repository:** %s/%s", pr.Owner, pr.Repo))
		if pr.TicketID != "" {
			result = append(result, fmt.Sprintf("**Ticket:** %s", pr.TicketID))
		}
	} else {
		result = append(result, fmt.Sprintf("## Found %d PRs:", len(prs)))
		result = append(result, "")

		for i, pr := range prs {
			result = append(result, fmt.Sprintf("%d. **#%d** - %s", i+1, pr.Number, pr.Title))
			result = append(result, fmt.Sprintf("   - URL: %s", pr.URL))
			result = append(result, fmt.Sprintf("   - State: %s | Branch: %s", pr.State, pr.HeadRefName))
			result = append(result, fmt.Sprintf("   - Repository: %s/%s", pr.Owner, pr.Repo))
			if pr.TicketID != "" {
				result = append(result, fmt.Sprintf("   - Ticket: %s", pr.TicketID))
			}
			result = append(result, "")
		}
	}

	return strings.Join(result, "\n")
}

// findBestPR finds the best/most relevant PR from a list
func findBestPR(prs []github.PR) github.PR {
	if len(prs) == 0 {
		return github.PR{}
	}

	// If only one PR, return it
	if len(prs) == 1 {
		return prs[0]
	}

	// Find the best PR based on priority:
	// 1. Open PRs (not closed/merged)
	// 2. Most recent PRs (highest number)
	// 3. PRs from main repositories

	var bestPR github.PR
	var bestScore int

	for _, pr := range prs {
		score := 0

		// Prioritize open PRs
		if pr.State == "OPEN" {
			score += 100
		}

		// Prioritize higher PR numbers (more recent)
		score += pr.Number

		// Prioritize main repositories (orbit, rasch-stack)
		if pr.Owner == "RingierAG" && (pr.Repo == "orbit" || pr.Repo == "rasch-stack") {
			score += 50
		}

		if score > bestScore {
			bestScore = score
			bestPR = pr
		}
	}

	return bestPR
}

// formatPRsJSON formats PRs in JSON format
func formatPRsJSON(prs []github.PR) string {
	// This would use the unified output system
	// For now, return a simple JSON representation
	return fmt.Sprintf(`{"prs": %d, "count": %d}`, len(prs), len(prs))
}
