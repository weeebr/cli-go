package git

import (
	"fmt"
	"sort"
	"strings"
)

// FormatRepoStats formats repository statistics for display
func FormatRepoStats(stats *RepoStats) string {
	var result strings.Builder

	result.WriteString("## Repository Statistics\n\n")
	result.WriteString(fmt.Sprintf("**Total Files:** %d\n", stats.TotalFiles))
	result.WriteString(fmt.Sprintf("**Total Lines:** %d\n", stats.TotalLines))
	result.WriteString(fmt.Sprintf("**Branches:** %d\n", stats.Branches))
	result.WriteString(fmt.Sprintf("**Commits:** %d\n", stats.Commits))
	result.WriteString(fmt.Sprintf("**Authors:** %d\n", stats.Authors))
	result.WriteString(fmt.Sprintf("**Lines Added:** %d\n", stats.LinesAdded))
	result.WriteString(fmt.Sprintf("**Lines Deleted:** %d\n", stats.LinesDeleted))

	if len(stats.Tickets) > 0 {
		result.WriteString(fmt.Sprintf("**Tickets:** %s\n", strings.Join(stats.Tickets, ", ")))
	}

	result.WriteString("\n## File Extensions\n\n")

	// Sort extensions by line count
	sort.Slice(stats.Extensions, func(i, j int) bool {
		return stats.Extensions[i].Lines > stats.Extensions[j].Lines
	})

	for _, ext := range stats.Extensions {
		result.WriteString(fmt.Sprintf("- **%s:** %d files, %d lines (%.1f%%)\n",
			ext.Extension, ext.Files, ext.Lines, float64(ext.Lines)/float64(stats.TotalLines)*100))
	}

	return result.String()
}

// FormatCommitStats formats commit statistics for display
func FormatCommitStats(stats *CommitStats) string {
	var result strings.Builder

	result.WriteString("## Commit Statistics\n\n")
	result.WriteString(fmt.Sprintf("**Total Commits:** %d\n", stats.TotalCommits))
	result.WriteString(fmt.Sprintf("**Total Lines Added:** %d\n", stats.TotalLinesAdded))
	result.WriteString(fmt.Sprintf("**Total Lines Deleted:** %d\n", stats.TotalLinesDeleted))

	result.WriteString("\n## Per Repository\n\n")

	// Sort repositories by commits
	type repoData struct {
		name  string
		stats RepoStats
	}
	var repos []repoData
	for name, repoStats := range stats.PerRepo {
		repos = append(repos, repoData{name, repoStats})
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].stats.Commits > repos[j].stats.Commits
	})

	for _, repo := range repos {
		result.WriteString(fmt.Sprintf("### %s\n", repo.name))
		result.WriteString(fmt.Sprintf("- Commits: %d\n", repo.stats.Commits))
		result.WriteString(fmt.Sprintf("- Lines Added: %d\n", repo.stats.LinesAdded))
		result.WriteString(fmt.Sprintf("- Lines Deleted: %d\n", repo.stats.LinesDeleted))
		result.WriteString(fmt.Sprintf("- Authors: %d\n", repo.stats.Authors))
		result.WriteString("\n")
	}

	result.WriteString("## Per Day\n\n")

	// Sort days by date
	var days []string
	for day := range stats.PerDay {
		days = append(days, day)
	}
	sort.Strings(days)

	for _, day := range days {
		dayStats := stats.PerDay[day]
		result.WriteString(fmt.Sprintf("### %s\n", day))
		result.WriteString(fmt.Sprintf("- Commits: %d\n", dayStats.Commits))
		result.WriteString(fmt.Sprintf("- Lines Added: %d\n", dayStats.LinesAdded))
		result.WriteString(fmt.Sprintf("- Lines Deleted: %d\n", dayStats.LinesDeleted))
		result.WriteString("\n")
	}

	if len(stats.TicketFrequency) > 0 {
		result.WriteString("## Ticket Frequency\n\n")

		// Sort tickets by frequency
		type ticketData struct {
			ticket string
			count  int
		}
		var tickets []ticketData
		for ticket, count := range stats.TicketFrequency {
			tickets = append(tickets, ticketData{ticket, count})
		}

		sort.Slice(tickets, func(i, j int) bool {
			return tickets[i].count > tickets[j].count
		})

		for _, ticket := range tickets {
			result.WriteString(fmt.Sprintf("- **%s:** %d commits\n", ticket.ticket, ticket.count))
		}
	}

	return result.String()
}

// FormatFileStats formats file statistics for display
func FormatFileStats(stats []FileStats) string {
	var result strings.Builder

	result.WriteString("## File Statistics\n\n")

	// Sort by line count
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Lines > stats[j].Lines
	})

	for _, stat := range stats {
		result.WriteString(fmt.Sprintf("- **%s:** %d files, %d lines (%.1f%%)\n",
			stat.Extension, stat.Files, stat.Lines, float64(stat.Lines)/float64(getTotalLines(stats))*100))
	}

	return result.String()
}

// FormatDayStats formats day statistics for display
func FormatDayStats(day string, stats DayStats) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("## %s\n\n", day))
	result.WriteString(fmt.Sprintf("**Commits:** %d\n", stats.Commits))
	result.WriteString(fmt.Sprintf("**Lines Added:** %d\n", stats.LinesAdded))
	result.WriteString(fmt.Sprintf("**Lines Deleted:** %d\n", stats.LinesDeleted))

	return result.String()
}

// getTotalLines calculates total lines from file stats
func getTotalLines(stats []FileStats) int {
	total := 0
	for _, stat := range stats {
		total += stat.Lines
	}
	return total
}

// FormatTicketFrequency formats ticket frequency for display
func FormatTicketFrequency(frequency map[string]int) string {
	var result strings.Builder

	result.WriteString("## Ticket Frequency\n\n")

	// Sort tickets by frequency
	type ticketData struct {
		ticket string
		count  int
	}
	var tickets []ticketData
	for ticket, count := range frequency {
		tickets = append(tickets, ticketData{ticket, count})
	}

	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].count > tickets[j].count
	})

	for _, ticket := range tickets {
		result.WriteString(fmt.Sprintf("- **%s:** %d commits\n", ticket.ticket, ticket.count))
	}

	return result.String()
}

// FormatSummary formats a summary of all statistics
func FormatSummary(repoStats *RepoStats, commitStats *CommitStats) string {
	var result strings.Builder

	result.WriteString("# Git Statistics Summary\n\n")

	// Repository summary
	result.WriteString("## Repository Overview\n")
	result.WriteString(fmt.Sprintf("- **Total Files:** %d\n", repoStats.TotalFiles))
	result.WriteString(fmt.Sprintf("- **Total Lines:** %d\n", repoStats.TotalLines))
	result.WriteString(fmt.Sprintf("- **Branches:** %d\n", repoStats.Branches))
	result.WriteString(fmt.Sprintf("- **Authors:** %d\n", repoStats.Authors))

	// Commit summary
	result.WriteString("\n## Commit Overview\n")
	result.WriteString(fmt.Sprintf("- **Total Commits:** %d\n", commitStats.TotalCommits))
	result.WriteString(fmt.Sprintf("- **Lines Added:** %d\n", commitStats.TotalLinesAdded))
	result.WriteString(fmt.Sprintf("- **Lines Deleted:** %d\n", commitStats.TotalLinesDeleted))

	// Top extensions
	result.WriteString("\n## Top File Extensions\n")
	sort.Slice(repoStats.Extensions, func(i, j int) bool {
		return repoStats.Extensions[i].Lines > repoStats.Extensions[j].Lines
	})

	for i, ext := range repoStats.Extensions {
		if i >= 5 { // Show top 5
			break
		}
		result.WriteString(fmt.Sprintf("- **%s:** %d files, %d lines\n",
			ext.Extension, ext.Files, ext.Lines))
	}

	return result.String()
}
