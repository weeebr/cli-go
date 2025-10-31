package main

// DESCRIPTION: show user activity across repositories

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"cli-go/_internal/git"
	"cli-go/_internal/github"
	"cli-go/_internal/io"
)

type Config struct {
	Compact bool
	Single  string
	Main    bool
	All     bool
	JSON    bool
}

func main() {

	cfg := parseFlags()

	// Get repositories to process
	// Default to main repos if no flags provided
	defaultMode := "main"
	if cfg.Single != "" || cfg.Main || cfg.All {
		defaultMode = "pwd" // Use explicit mode if flags provided
	}
	repoPaths, err := git.GetReposToProcess(cfg.Single, cfg.Main, cfg.All, defaultMode)
	ai.ExitIf(err, "failed to get repositories to process")

	// Load config for user info
	configData, err := config.LoadConfig()
	ai.ExitIf(err, "failed to load config")

	// Get branches using git library
	var allBranches []git.Branch
	for _, repoPath := range repoPaths {
		branches, err := git.GetUserBranches(repoPath, configData.Ringier.DefaultUser)
		if err != nil {
			io.LogWarning("Failed to get branches from %s: %v", repoPath, err)
			continue
		}
		allBranches = append(allBranches, branches...)
	}

	// Get PRs using github library
	var allPRs []github.PR
	for _, repoConfig := range configData.Repositories {
		prs, err := github.GetUserOpenPRs(configData.Ringier.DefaultUser, repoConfig.Owner, repoConfig.Repo)
		if err != nil {
			io.LogWarning("Failed to get PRs from %s/%s: %v", repoConfig.Owner, repoConfig.Repo, err)
			continue
		}
		allPRs = append(allPRs, prs...)
	}

	// Format output
	activity := map[string]interface{}{
		"branches": allBranches,
		"prs":      allPRs,
	}

	if cfg.JSON {
		io.DirectOutput(activity, *clip, *file, true)
	} else {
		outputDefault(activity)
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

	flag.Parse()

	// No default needed - GetReposToProcess handles defaults

	return config
}

func outputDefault(activity map[string]interface{}) {
	if branches, exists := activity["branches"]; exists {
		if branchList, ok := branches.([]git.Branch); ok {
			if len(branchList) > 0 {
				fmt.Printf("🌿 Your branches:\n")
				for _, branch := range branchList {
					fmt.Printf("  • %s (%s)\n", branch.Name, branch.RepoPath)
				}
			} else {
				fmt.Printf("🌿 No branches found\n")
			}
		}
	}

	if prs, exists := activity["prs"]; exists {
		if prList, ok := prs.([]github.PR); ok {
			if len(prList) > 0 {
				fmt.Printf("\n📋 Open PRs:\n")
				for _, pr := range prList {
					fmt.Printf("  • %s (#%d) - %s\n", pr.Title, pr.Number, pr.Repo)
				}
			} else {
				fmt.Printf("\n📋 No open PRs found\n")
			}
		}
	}
}
