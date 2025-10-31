package main

// DESCRIPTION: show file / LOC stats of repo

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/flags"
	"cli-go/_internal/git"
	"cli-go/_internal/io"
	"sort"
	"strings"
)

type Config struct {
	Compact bool
	Single  string
	Main    bool
	All     bool
	JSON    bool
}

func main() {
	config := parseFlags()

	// Check for help command
	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		io.LogInfo("gstats - Show file/LOC stats of repository")
		io.LogInfo("Analyzes git-tracked files and provides statistics by extension")
		io.LogInfo("Supports --single <path> (specific repo), --main (main repos), --all (all repos)")
		io.LogInfo("Default: current repository")
		io.LogInfo("Output: Formatted display (default) or JSON with --json flag")
		return
	}

	// Determine which repositories to process
	repoPaths, err := git.GetReposToProcess(config.Single, config.Main, config.All, "pwd")
	ai.ExitIf(err, "failed to get repositories to process")

	var result interface{}
	if len(repoPaths) == 1 {
		// Single repo - use existing behavior
		stats, err := git.GetRepoStats(repoPaths[0])
		ai.ExitIf(err, "failed to get repo stats")
		result = stats
	} else {
		// Multiple repos - aggregate stats
		// For multiple repos, we'll need to implement this
		// For now, just use the first repo
		stats, err := git.GetRepoStats(repoPaths[0])
		ai.ExitIf(err, "failed to get repo stats")
		repoStats := []git.RepoStats{*stats}
		// No error to check here - this is just assignment
		result = repoStats
	}

	if config.JSON {
		io.DirectOutput(result, *clip, *file, true)
	} else {
		// Formatted output for better display
		formattedOutput := formatStatsOutput(result)
		io.DirectOutput(formattedOutput, *clip, *file, false)
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
	flag.BoolVar(&config.JSON, "json", false, "Output in JSON format (default: formatted)")

	flags.ReorderAndParse()

	// No default needed - GetReposToProcess handles defaults

	return config
}

// formatStatsOutput creates a nicely formatted output similar to legacy CLI tools
func formatStatsOutput(result interface{}) string {
	var output strings.Builder

	// Handle single repo stats
	if stats, ok := result.(*git.RepoStats); ok {
		output.WriteString(io.FormatWithEmoji("Repository Statistics", "cache") + "\n")
		output.WriteString("=" + strings.Repeat("=", 30) + "\n\n")

		// Summary stats
		output.WriteString(fmt.Sprintf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Total Files: %d", stats.TotalFiles), "file")))
		output.WriteString(fmt.Sprintf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Total Lines: %d", stats.TotalLines), "file")))
		output.WriteString(fmt.Sprintf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Branches: %d", stats.Branches), "branches")))
		output.WriteString(fmt.Sprintf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Commits: %d", stats.Commits), "commits")))
		output.WriteString(fmt.Sprintf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Authors: %d", stats.Authors), "user")))
		output.WriteString(fmt.Sprintf("➕ Lines Added: %d\n", stats.LinesAdded))
		output.WriteString(fmt.Sprintf("➖ Lines Deleted: %d\n\n", stats.LinesDeleted))

		// File extensions breakdown
		if len(stats.Extensions) > 0 {
			output.WriteString("📋 Files by Extension:\n")
			output.WriteString("-" + strings.Repeat("-", 25) + "\n")

			// Sort extensions by line count (descending)
			sort.Slice(stats.Extensions, func(i, j int) bool {
				return stats.Extensions[i].Lines > stats.Extensions[j].Lines
			})

			for _, ext := range stats.Extensions {
				output.WriteString(fmt.Sprintf("  %-8s %4d files  %6d lines (%2d%%)\n",
					ext.Extension, ext.Files, ext.Lines, ext.Percent))
			}
		}

		// Tickets if any
		if len(stats.Tickets) > 0 {
			output.WriteString(fmt.Sprintf("\n%s\n", io.FormatWithEmoji(fmt.Sprintf("Recent Tickets: %s", strings.Join(stats.Tickets, ", ")), "ticket")))
		}
	} else if repoStats, ok := result.([]git.RepoStats); ok {
		// Handle multiple repo stats
		output.WriteString(io.FormatWithEmoji("Multi-Repository Statistics", "cache") + "\n")
		output.WriteString("=" + strings.Repeat("=", 35) + "\n\n")

		totalFiles := 0
		totalLines := 0
		totalBranches := 0
		totalCommits := 0
		totalAuthors := 0

		for i, stats := range repoStats {
			output.WriteString(fmt.Sprintf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Repository %d:", i+1), "file")))
			output.WriteString(fmt.Sprintf("  %s | %s | %s\n",
				io.FormatWithEmoji(fmt.Sprintf("Files: %d", stats.TotalFiles), "file"),
				io.FormatWithEmoji(fmt.Sprintf("Lines: %d", stats.TotalLines), "file"),
				io.FormatWithEmoji(fmt.Sprintf("Branches: %d", stats.Branches), "branches")))
			output.WriteString(fmt.Sprintf("  %s | %s\n\n",
				io.FormatWithEmoji(fmt.Sprintf("Commits: %d", stats.Commits), "commits"),
				io.FormatWithEmoji(fmt.Sprintf("Authors: %d", stats.Authors), "user")))

			totalFiles += stats.TotalFiles
			totalLines += stats.TotalLines
			totalBranches += stats.Branches
			totalCommits += stats.Commits
			totalAuthors += stats.Authors
		}

		output.WriteString(io.FormatWithEmoji("Totals:", "cache") + "\n")
		output.WriteString(fmt.Sprintf("  %s | %s | %s\n",
			io.FormatWithEmoji(fmt.Sprintf("Files: %d", totalFiles), "file"),
			io.FormatWithEmoji(fmt.Sprintf("Lines: %d", totalLines), "file"),
			io.FormatWithEmoji(fmt.Sprintf("Branches: %d", totalBranches), "branches")))
		output.WriteString(fmt.Sprintf("  %s | %s\n",
			io.FormatWithEmoji(fmt.Sprintf("Commits: %d", totalCommits), "commits"),
			io.FormatWithEmoji(fmt.Sprintf("Authors: %d", totalAuthors), "user")))
	}

	return output.String()
}
