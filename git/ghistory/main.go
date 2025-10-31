package main

// DESCRIPTION: show commit history across repositories

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/flags"
	"cli-go/_internal/git"
	"cli-go/_internal/io"
	"os"
	"time"
)

type Config struct {
	Compact bool
	Single  string
	Main    bool
	All     bool
	Days    int
	Author  string
	JSON    bool
}

func main() {
	config := parseFlags()

	// Check for help command
	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		io.LogInfo("ghistory - Show commit history across repositories")
		io.LogInfo("Shows commit history with filtering by days and author")
		io.LogInfo("Supports --single <path> (specific repo), --main (main repos), --all (all repos)")
		io.LogInfo("Default: all repositories")
		io.LogInfo("Flags: --days N (days to look back), --author 'name' (filter by author)")
		io.LogInfo("Output: JSON with commit history grouped by day")
		return
	}

	// Determine which repositories to process
	repoPaths, err := git.GetReposToProcess(config.Single, config.Main, config.All, "all")
	ai.ExitIf(err, "failed to get git history")

	// Calculate since time
	since := time.Now().AddDate(0, 0, -config.Days)

	// Get commits from all repositories
	var allCommits []git.Commit
	for _, repoPath := range repoPaths {
		commits, err := git.GetCommits(repoPath, since, config.Author)
		if err != nil {
			// Log warning but continue with other repos
			io.LogWarning("Failed to get commits from %s: %v", repoPath, err)
			continue
		}
		allCommits = append(allCommits, commits...)
	}

	if len(allCommits) == 0 {
		io.LogError("No commits found")
		os.Exit(1)
	}

	// Format and output results
	if config.JSON {
		io.DirectOutput(allCommits, *clip, *file, true)
	} else {
		// Format as markdown for human-readable output
		output := formatHistoryMarkdown(allCommits)
		io.LogInfo("%s", output)
	}
}

func formatHistoryMarkdown(commits []git.Commit) string {
	if len(commits) == 0 {
		return "No commits found."
	}

	// Group commits by day
	commitsByDay := make(map[string][]git.Commit)
	for _, commit := range commits {
		day := commit.Date.Format("2006-01-02")
		commitsByDay[day] = append(commitsByDay[day], commit)
	}

	// Sort days
	var days []string
	for day := range commitsByDay {
		days = append(days, day)
	}

	// Sort days in descending order (most recent first)
	for i := 0; i < len(days)-1; i++ {
		for j := i + 1; j < len(days); j++ {
			if days[i] < days[j] {
				days[i], days[j] = days[j], days[i]
			}
		}
	}

	var result []string
	result = append(result, "## Commit History")
	result = append(result, "")

	for _, day := range days {
		dayCommits := commitsByDay[day]
		result = append(result, "### "+day+" ("+string(rune(len(dayCommits)))+" commits)")
		result = append(result, "")

		for _, commit := range dayCommits {
			result = append(result, "- **"+commit.Hash[:8]+"** - "+commit.Message)
			result = append(result, fmt.Sprintf("  - %s", io.FormatWithEmoji(fmt.Sprintf("Author: %s", commit.Author), "user")))
			result = append(result, "  - Repository: "+commit.Repository)
			result = append(result, "  - Branch: "+commit.Branch)
			if len(commit.TicketIDs) > 0 {
				result = append(result, "  - Tickets: "+joinStrings(commit.TicketIDs, ", "))
			}
			if commit.LinesAdded > 0 || commit.LinesDeleted > 0 {
				result = append(result, "  - Changes: +"+string(rune(commit.LinesAdded))+" -"+string(rune(commit.LinesDeleted)))
			}
			result = append(result, "")
		}
	}

	return joinStrings(result, "\n")
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

var (
	clip = flag.Bool("clip", false, "Copy to clipboard")
	file = flag.String("file", "", "Write to file")
)

func parseFlags() Config {
	config := Config{
		Days: 7, // Default to 7 days
	}

	flag.BoolVar(&config.Compact, "compact", false, "Use compact JSON format")
	flag.StringVar(&config.Single, "single", "", "Operate on specific repository path")
	flag.BoolVar(&config.Main, "main", false, "Operate on main repositories (orbit + rasch-stack)")
	flag.BoolVar(&config.All, "all", false, "Operate on all repositories from config")
	flag.IntVar(&config.Days, "days", 7, "Number of days to look back")
	flag.StringVar(&config.Author, "author", "", "Filter by author name")
	flag.BoolVar(&config.JSON, "json", false, "Output in JSON format")

	flags.ReorderAndParse()

	// No default needed - GetReposToProcess handles defaults

	return config
}
