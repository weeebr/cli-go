package main

// DESCRIPTION: show current files vs branched off

import (
	"flag"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/flags"
	"cli-go/_internal/git"
	"cli-go/_internal/io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Compact bool
	JSON    bool
	Test    bool
}

func main() {
	config := parseFlags()

	// Check for help command
	args := flag.Args()
	if len(args) > 0 && args[0] == "help" {
		io.LogInfo("gaff - Show changed files vs fork-point")
		io.LogInfo("Shows files that have changed since the fork-point branch")
		io.LogInfo("Opens selected file in diff mode with editor")
		io.LogInfo("Output: {\"files\": [\"file1.go\", \"file2.ts\"], \"base_branch\": \"origin/main\", \"selected\": \"file1.go\"}")
		return
	}

	// Check if a specific file was provided as argument (skip flags)
	var selectedFile string
	fileArgs := flag.Args()
	if len(fileArgs) > 0 && !strings.HasPrefix(fileArgs[0], "-") {
		selectedFile = fileArgs[0]
	}

	forkInfo, err := git.GetChangedFilesWithSelection(selectedFile, config.Test)
	ai.ExitIf(err, "failed to get changed files")

	// If no files changed, show message and exit
	if len(forkInfo.Files) == 0 {
		io.LogInfo("no changes since fork-point")
		return
	}

	// If a specific file was selected, open it in diff mode
	if selectedFile != "" {
		if err := openFileInDiffMode(selectedFile, forkInfo.BaseBranch); err != nil {
			ai.ExitIf(err, "failed to open file in diff mode")
		}
		return
	}

	// If no specific file, show output for interactive selection
	if config.JSON {
		io.DirectOutput(forkInfo, *clip, *file, true)
	} else {
		outputDefault(forkInfo)
	}
}

var (
	clip = flag.Bool("clip", false, "Copy to clipboard")
	file = flag.String("file", "", "Write to file")
)

func parseFlags() Config {
	config := Config{}

	flag.BoolVar(&config.Compact, "compact", false, "Use compact JSON format")
	flag.BoolVar(&config.JSON, "json", false, "Output in JSON format")
	flag.BoolVar(&config.Test, "test", false, "Test mode - pre-select first changed file")

	flags.ReorderAndParse()

	return config
}

func openFileInDiffMode(filePath, baseBranch string) error {
	// Get repository root
	repoRoot, err := git.GetGitRoot()
	if err != nil {
		return fmt.Errorf("failed to get repo root: %v", err)
	}

	// Get absolute path to the file
	var absPath string
	if filepath.IsAbs(filePath) {
		absPath = filePath
	} else {
		absPath = filepath.Join(repoRoot, filePath)
	}

	// Get relative path from repo root
	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %v", err)
	}

	// Get current branch
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Create temporary directory for fork-point version
	tmpDir := fmt.Sprintf("%s/gaff-%s-%s-%d",
		os.TempDir(),
		filepath.Base(repoRoot),
		currentBranch,
		os.Getpid())

	tmpFile := filepath.Join(tmpDir, relPath)

	// Create directory structure
	if err := os.MkdirAll(filepath.Dir(tmpFile), 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}

	// Get fork-point version of the file
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", baseBranch, relPath))
	output, err := cmd.Output()
	if err != nil {
		// If file doesn't exist in fork-point, create empty file
		output = []byte{}
	}

	// Write fork-point version to temp file
	if err := os.WriteFile(tmpFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "code" // Default to VS Code
	}

	// Parse editor command (handle cases like "cursor -w")
	editorParts := strings.Fields(editor)
	editorCmd := editorParts[0]
	editorArgs := editorParts[1:]
	
	// Remove -w flag as it's incompatible with --diff
	var filteredArgs []string
	for _, arg := range editorArgs {
		if arg != "-w" && arg != "--wait" {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// Open editor in diff mode
	// Cursor uses -d flag (without -r as it may conflict)
	args := append(filteredArgs, "-d", tmpFile, absPath)
	
	// Debug output
	fmt.Fprintf(os.Stderr, "Opening diff: %s %v\n", editorCmd, args)
	
	cmd = exec.Command(editorCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open editor: %v", err)
	}
	
	fmt.Fprintf(os.Stderr, "Diff opened successfully\n")

	// Clean up temp directory after a delay
	go func() {
		// Wait longer for the editor to fully load the files
		time.Sleep(30 * time.Second)
		os.RemoveAll(tmpDir)
	}()

	return nil
}

func outputDefault(forkInfo interface{}) {
	if resultMap, ok := forkInfo.(map[string]interface{}); ok {
		if files, exists := resultMap["files"]; exists {
			if fileList, ok := files.([]string); ok {
				if len(fileList) > 0 {
					fmt.Printf("%s\n", io.FormatWithEmoji("Changed files since fork-point:", "file"))
					for _, file := range fileList {
						fmt.Printf("  • %s\n", file)
					}
				} else {
					fmt.Printf("%s\n", io.FormatWithEmoji("No changes since fork-point", "file"))
				}
			}
		}
		if baseBranch, exists := resultMap["base_branch"]; exists {
			fmt.Printf("%s\n", io.FormatWithEmoji(fmt.Sprintf("Base branch: %s", baseBranch), "branches"))
		}
	}
}
