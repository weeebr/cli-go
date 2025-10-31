package git

import (
	"bufio"
	"fmt"
	"cli-go/_internal/sys"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// validExts defines which file extensions are considered source files
var validExts = map[string]bool{
	"go": true, "js": true, "ts": true, "py": true, "java": true,
	"cpp": true, "c": true, "h": true, "hpp": true, "cs": true,
	"php": true, "rb": true, "swift": true, "kt": true, "scala": true,
	"rs": true, "dart": true, "vue": true, "jsx": true, "tsx": true,
	"svelte": true, "html": true, "css": true, "scss": true,
	"sass": true, "less": true, "xml": true, "json": true,
	"yaml": true, "yml": true, "toml": true, "ini": true,
	"cfg": true, "conf": true, "sh": true, "bash": true,
	"zsh": true, "fish": true, "ps1": true, "bat": true,
	"cmd": true, "sql": true, "r": true, "m": true,
	"pl": true, "pm": true, "t": true, "pod": true,
	"tex": true, "bib": true, "sty": true, "cls": true,
	"dtx": true, "ins": true, "lua": true, "vim": true,
	"el": true, "hs": true, "lhs": true, "ml": true,
	"mli": true, "fs": true, "fsi": true, "fsx": true,
	"f90": true, "f95": true, "f03": true, "f08": true,
	"f": true, "for": true, "f77": true, "f66": true,
}

// FileStats represents statistics for a file extension
type FileStats struct {
	Extension string `json:"extension"`
	Files     int    `json:"files"`
	Lines     int    `json:"lines"`
	Percent   int    `json:"percent"`
	AvgLines  int    `json:"avg_lines"`
}

// RepoStats represents statistics for a single repository
type RepoStats struct {
	Extensions   []FileStats `json:"extensions"`
	TotalFiles   int         `json:"total_files"`
	TotalLines   int         `json:"total_lines"`
	Branches     int         `json:"branches"`
	Commits      int         `json:"commits"`
	Authors      int         `json:"authors"`
	LinesAdded   int         `json:"linesAdded"`
	LinesDeleted int         `json:"linesDeleted"`
	Tickets      []string    `json:"tickets"`
}

// CommitStats represents commit statistics
type CommitStats struct {
	TotalCommits      int                  `json:"totalCommits"`
	TotalLinesAdded   int                  `json:"totalLinesAdded"`
	TotalLinesDeleted int                  `json:"totalLinesDeleted"`
	PerRepo           map[string]RepoStats `json:"perRepo"`
	PerDay            map[string]DayStats  `json:"perDay"`
	TicketFrequency   map[string]int       `json:"ticketFrequency"`
}

// DayStats represents statistics for a single day
type DayStats struct {
	Commits      int      `json:"commits"`
	LinesAdded   int      `json:"linesAdded"`
	LinesDeleted int      `json:"linesDeleted"`
	Authors      []string `json:"authors"`
}

// findSourceFiles finds all source files in the repository using native Go
func findSourceFiles(repoPath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if validExts[ext] {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// GetRepoStats calculates statistics for a repository
func GetRepoStats(repoPath string) (*RepoStats, error) {
	// Get all files in the repository using native Go
	filePaths, err := findSourceFiles(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find files: %v", err)
	}

	// Convert to string for compatibility with existing code
	files := strings.Join(filePaths, "\n")

	extStats := make(map[string]*FileStats)
	totalFiles := 0
	totalLines := 0

	// Process each file
	fileLines := strings.Split(files, "\n")
	for _, file := range fileLines {
		if file == "" {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file))
		if ext != "" {
			ext = ext[1:] // Remove the dot
		}

		if !validExts[ext] {
			continue
		}

		if extStats[ext] == nil {
			extStats[ext] = &FileStats{Extension: ext}
		}

		// Count lines in file
		// file already contains the full path from find command
		fullPath := file
		lines, err := CountLinesInFile(fullPath)
		if err != nil {
			continue // Skip files that can't be read
		}

		extStats[ext].Files++
		extStats[ext].Lines += lines
		totalFiles++
		totalLines += lines
	}

	// Convert to slice and calculate percentages
	var extensions []FileStats
	for _, stats := range extStats {
		if totalLines > 0 {
			stats.Percent = int((float64(stats.Lines) / float64(totalLines)) * 100)
		}
		if stats.Files > 0 {
			stats.AvgLines = stats.Lines / stats.Files
		}
		extensions = append(extensions, *stats)
	}

	// Sort by lines (descending)
	sort.Slice(extensions, func(i, j int) bool {
		return extensions[i].Lines > extensions[j].Lines
	})

	// Get additional stats
	branches := GetBranchCount(repoPath)
	commits, err := GetCommitCount(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit count: %v", err)
	}
	authors, err := GetAuthorCount(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get author count: %v", err)
	}

	return &RepoStats{
		Extensions: extensions,
		TotalFiles: totalFiles,
		TotalLines: totalLines,
		Branches:   branches,
		Commits:    commits,
		Authors:    authors,
	}, nil
}

// CalculateCommitStats calculates statistics from commits
func CalculateCommitStats(commits []Commit) CommitStats {
	stats := CommitStats{
		PerRepo:         make(map[string]RepoStats),
		PerDay:          make(map[string]DayStats),
		TicketFrequency: make(map[string]int),
	}

	for _, commit := range commits {
		// Update totals
		stats.TotalCommits++
		stats.TotalLinesAdded += commit.LinesAdded
		stats.TotalLinesDeleted += commit.LinesDeleted

		// Update per-repo stats
		repoName := commit.Repository
		if repoStats, exists := stats.PerRepo[repoName]; exists {
			repoStats.Commits++
			repoStats.LinesAdded += commit.LinesAdded
			repoStats.LinesDeleted += commit.LinesDeleted
			repoStats.Tickets = append(repoStats.Tickets, commit.TicketIDs...)
			stats.PerRepo[repoName] = repoStats
		} else {
			stats.PerRepo[repoName] = RepoStats{
				Commits:      1,
				LinesAdded:   commit.LinesAdded,
				LinesDeleted: commit.LinesDeleted,
				Tickets:      commit.TicketIDs,
			}
		}

		// Update per-day stats
		day := commit.Date.Format("2006-01-02")
		if dayStats, exists := stats.PerDay[day]; exists {
			dayStats.Commits++
			dayStats.LinesAdded += commit.LinesAdded
			dayStats.LinesDeleted += commit.LinesDeleted
			stats.PerDay[day] = dayStats
		} else {
			stats.PerDay[day] = DayStats{
				Commits:      1,
				LinesAdded:   commit.LinesAdded,
				LinesDeleted: commit.LinesDeleted,
			}
		}

		// Update ticket frequency
		for _, ticket := range commit.TicketIDs {
			stats.TicketFrequency[ticket]++
		}
	}

	return stats
}

// CountLinesInFile counts the number of lines in a file
func CountLinesInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	return lines, scanner.Err()
}

// GetCommitCount returns the number of commits in the repository
func GetCommitCount(repoPath string) (int, error) {
	result := sys.RunCommand("git", "-C", repoPath, "rev-list", "--count", "HEAD")
	if result.Error != nil {
		return 0, result.Error
	}
	output := result.Stdout

	var count int
	fmt.Sscanf(strings.TrimSpace(output), "%d", &count)
	return count, nil
}

// GetAuthorCount returns the number of unique authors in the repository
func GetAuthorCount(repoPath string) (int, error) {
	// Get all authors
	result := sys.RunCommand("git", "-C", repoPath, "log", "--pretty=format:%an")
	if result.Error != nil {
		return 0, result.Error
	}

	// Count unique authors
	authors := strings.Split(result.Stdout, "\n")
	uniqueAuthors := make(map[string]bool)
	for _, author := range authors {
		author = strings.TrimSpace(author)
		if author != "" {
			uniqueAuthors[author] = true
		}
	}

	return len(uniqueAuthors), nil
}
