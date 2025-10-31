package git

import (
	"fmt"
	"cli-go/_internal/sys"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ForkPointInfo holds information about fork point detection
type ForkPointInfo struct {
	BaseBranch string   `json:"base_branch"`
	BaseCommit string   `json:"base_commit"`
	Files      []string `json:"files"`
}

// ForkPointDiffInfo holds information for diff operations
type ForkPointDiffInfo struct {
	BaseBranch string   `json:"base_branch"`
	BaseCommit string   `json:"base_commit"`
	Files      []string `json:"files"`
	Selected   string   `json:"selected,omitempty"`
	TempDir    string   `json:"temp_dir,omitempty"`
}

// GetForkPointBranch detects the most likely branched-from remote branch
func GetForkPointBranch() (string, error) {
	// First check if we have any remote branches
	result := sys.RunCommand("git", "for-each-ref", "--count=1", "refs/remotes/origin")
	if result.ExitCode != 0 || result.Stdout == "" {
		return "", fmt.Errorf("no remote branches found")
	}

	// Get current upstream to exclude it
	currentUpstream := ""
	upstreamResult := sys.RunCommand("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}")
	if upstreamResult.ExitCode == 0 && upstreamResult.Stdout != "" {
		currentUpstream = upstreamResult.Stdout
	}

	// Get recent remote branches (simplified approach)
	result = sys.RunCommand("git", "for-each-ref", "--sort=-committerdate", "--count=5", "--format=%(refname:short)", "refs/remotes/origin")
	if result.ExitCode != 0 {
		return "", fmt.Errorf("failed to get remote branches: %s", result.Stderr)
	}

	lines := strings.Split(result.Stdout, "\n")

	// Return the first valid remote branch (excluding current upstream)
	for _, line := range lines {
		if line == "" || line == "origin/HEAD" {
			continue
		}

		// Skip current upstream
		if currentUpstream != "" && line == currentUpstream {
			continue
		}

		return line, nil
	}

	// If we only have one remote branch and it's the current upstream, use it anyway
	if len(lines) > 0 && lines[0] != "" && lines[0] != "origin/HEAD" {
		return lines[0], nil
	}

	return "", fmt.Errorf("could not determine base branch")
}

// GetForkPointCommit gets the merge-base commit between HEAD and base branch
func GetForkPointCommit(baseBranch string) (string, error) {
	result := sys.RunCommand("git", "merge-base", "HEAD", baseBranch)
	if result.ExitCode != 0 {
		return "", fmt.Errorf("failed to get fork point: %s", result.Stderr)
	}
	return result.Stdout, nil
}

// GetChangedFilesSinceForkPoint returns files changed since fork point
func GetChangedFilesSinceForkPoint() (*ForkPointInfo, error) {
	baseBranch, err := GetForkPointBranch()
	if err != nil {
		return nil, err
	}

	baseCommit, err := GetForkPointCommit(baseBranch)
	if err != nil {
		return nil, err
	}

	// Get changed files
	result := sys.RunCommand("git", "diff", "--name-only", baseCommit)
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get changed files: %s", result.Stderr)
	}

	var files []string
	if result.Stdout != "" {
		files = strings.Split(result.Stdout, "\n")
		// Remove empty strings
		var cleanFiles []string
		for _, file := range files {
			if file != "" {
				cleanFiles = append(cleanFiles, file)
			}
		}
		files = cleanFiles
	}

	return &ForkPointInfo{
		BaseBranch: baseBranch,
		BaseCommit: baseCommit,
		Files:      files,
	}, nil
}

// GetChangedFilesWithSelection returns files changed since fork point with interactive selection
func GetChangedFilesWithSelection(selectedFile string, testMode bool) (*ForkPointDiffInfo, error) {
	baseBranch, err := GetForkPointBranch()
	if err != nil {
		return nil, err
	}

	baseCommit, err := GetForkPointCommit(baseBranch)
	if err != nil {
		return nil, err
	}

	// Get changed files
	result := sys.RunCommand("git", "diff", "--name-only", baseCommit)
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get changed files: %s", result.Stderr)
	}

	var files []string
	if result.Stdout != "" {
		fileList := strings.Split(result.Stdout, "\n")
		// Remove empty strings
		for _, file := range fileList {
			if file != "" {
				files = append(files, file)
			}
		}
	}

	// If no files changed, return early
	if len(files) == 0 {
		return &ForkPointDiffInfo{
			BaseBranch: baseBranch,
			BaseCommit: baseCommit,
			Files:      files,
		}, nil
	}

	// If no file selected, use fzf for interactive selection
	selected := selectedFile
	if selected == "" {
		if testMode {
			// Test mode: pre-select first file
			if len(files) > 0 {
				selected = files[0]
			} else {
				// No files to select, return all files
				return &ForkPointDiffInfo{
					BaseBranch: baseBranch,
					BaseCommit: baseCommit,
					Files:      files,
				}, nil
			}
		} else {
			// Check if fzf is available
			if _, err := exec.LookPath("fzf"); err == nil {
				// Use fzf for selection with pre-selection for non-interactive environments
				filesStr := strings.Join(files, "\n")
				// Pre-select first file for non-interactive environments
				preSelect := ""
				if len(files) > 0 {
					preSelect = fmt.Sprintf("--query=%s", files[0])
				}
				cmd := exec.Command("fzf", "--prompt=gaff> ", "--ansi", "--border", "--height=80%", preSelect)
				cmd.Stdin = strings.NewReader(filesStr)
				output, err := cmd.Output()
				if err == nil && len(output) > 0 {
					selected = strings.TrimSpace(string(output))
				} else {
					// fzf failed or was cancelled, return all files
					return &ForkPointDiffInfo{
						BaseBranch: baseBranch,
						BaseCommit: baseCommit,
						Files:      files,
					}, nil
				}
			} else {
				// Fallback: return all files
				return &ForkPointDiffInfo{
					BaseBranch: baseBranch,
					BaseCommit: baseCommit,
					Files:      files,
				}, nil
			}
		}
	}

	// If no file was selected (fzf failed or no fzf), return all files
	if selected == "" {
		return &ForkPointDiffInfo{
			BaseBranch: baseBranch,
			BaseCommit: baseCommit,
			Files:      files,
		}, nil
	}

	// Get repo root
	repoRoot, err := GetGitRoot()
	if err != nil {
		return nil, err
	}

	// Create temp directory for fork-point snapshots
	tempDir := fmt.Sprintf("/tmp/gaff-%s-%s-%d",
		strings.ReplaceAll(repoRoot, "/", "_"),
		baseBranch,
		os.Getpid())

	// Create temp file for selected file
	relPath := selected
	if strings.HasPrefix(selected, repoRoot) {
		relPath = strings.TrimPrefix(selected, repoRoot+"/")
	}

	tempFile := fmt.Sprintf("%s/%s", tempDir, relPath)

	// Create temp directory
	os.MkdirAll(filepath.Join(tempDir, filepath.Dir(relPath)), 0755)

	// Create fork-point snapshot
	snapshotResult := sys.RunCommand("git", "show", fmt.Sprintf("%s:%s", baseBranch, relPath))
	if snapshotResult.ExitCode == 0 {
		// Write snapshot to temp file
		os.WriteFile(tempFile, []byte(snapshotResult.Stdout), 0644)
	} else {
		// Create empty file if snapshot doesn't exist
		os.Create(tempFile)
	}

	// Editor launching logic moved to gaff tool

	return &ForkPointDiffInfo{
		BaseBranch: baseBranch,
		BaseCommit: baseCommit,
		Files:      files,
		Selected:   selected,
		TempDir:    tempDir,
	}, nil
}

// Test comment
