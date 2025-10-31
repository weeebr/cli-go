package git

import (
	"fmt"
	"cli-go/_internal/custom"
	"cli-go/_internal/sys"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Commit represents a Git commit with metadata
type Commit struct {
	Hash         string    `json:"hash"`
	Message      string    `json:"message"`
	Author       string    `json:"author"`
	Date         time.Time `json:"date"`
	Repository   string    `json:"repository"`
	Branch       string    `json:"branch"`
	TicketIDs    []string  `json:"ticketIds"`
	FilesAdded   int       `json:"filesAdded"`
	FilesDeleted int       `json:"filesDeleted"`
	LinesAdded   int       `json:"linesAdded"`
	LinesDeleted int       `json:"linesDeleted"`
}

// GetCommits gets commits from a repository since a specific time
func GetCommits(repoPath string, since time.Time, author string) ([]Commit, error) {
	// Check if it's a git repository
	if !IsGitRepo() {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Build git log command
	args := []string{"log", "--since", since.Format("2006-01-02"), "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso"}
	if author != "" {
		args = append(args, "--author", author)
	}

	result := sys.RunCommand("git", args...)
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get commits: %s", result.Stderr)
	}

	var commits []Commit
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		hash := parts[0]
		authorName := parts[1]
		_ = parts[2] // authorEmail - not used
		dateStr := parts[3]
		message := parts[4]

		date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			// Try alternative format
			date, err = time.Parse("2006-01-02T15:04:05-07:00", dateStr)
			if err != nil {
				continue
			}
		}

		// Get branch name
		branch, _ := GetCurrentBranch()

		// Extract ticket IDs from message
		tickets := custom.ExtractTicketsFromMessage(message)

		// Get diff stats
		added, deleted, _ := GetDiffStats(repoPath, hash)

		commit := Commit{
			Hash:         hash,
			Message:      message,
			Author:       authorName,
			Date:         date,
			Repository:   repoPath,
			Branch:       branch,
			TicketIDs:    tickets,
			LinesAdded:   added,
			LinesDeleted: deleted,
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// GetDiffStats gets the diff statistics for a commit
func GetDiffStats(repoPath, commitHash string) (added, deleted int, err error) {
	if !IsGitRepo() {
		return 0, 0, fmt.Errorf("not a git repository: %s", repoPath)
	}

	result := sys.RunCommand("git", "show", "--stat", "--format=", commitHash)
	if result.ExitCode != 0 {
		return 0, 0, fmt.Errorf("failed to get diff stats: %s", result.Stderr)
	}

	// Parse the diff stats output
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if strings.Contains(line, "insertion") || strings.Contains(line, "deletion") {
			// Parse lines like " 2 files changed, 3 insertions(+), 1 deletion(-)"
			re := regexp.MustCompile(`(\d+) insertion.*?(\d+) deletion`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				if ins, err := strconv.Atoi(matches[1]); err == nil {
					added += ins
				}
				if del, err := strconv.Atoi(matches[2]); err == nil {
					deleted += del
				}
			}
		}
	}

	return added, deleted, nil
}
