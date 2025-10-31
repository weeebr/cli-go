package main

// DESCRIPTION: check what's behind an alias

import (
	"flag"
	"fmt"
	"cli-go/_internal/config"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CheckAliasResult struct {
	Command string `json:"command"`
	Type    string `json:"type"` // "alias", "function", "git_alias", "binary", "not_found"
	Value   string `json:"value,omitempty"`
	File    string `json:"file,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
}

func main() {

	var (
		clip = flag.Bool("clip", false, "Copy to clipboard")
		file = flag.String("file", "", "Write to file")
		json = flag.Bool("json", false, "Output in JSON format")
	)
	flags.ReorderAndParse()

	args := flag.Args()
	if len(args) == 0 {
		result := CheckAliasResult{
			Error:   "command required",
			Message: "Usage: check_alias <command>",
		}
		if *json {
			// Use direct output
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		os.Exit(1)
	}

	command := args[0]

	// Check if command exists
	path, err := exec.LookPath(command)
	if err != nil {
		result := CheckAliasResult{
			Command: command,
			Type:    "not_found",
			Message: "Command not found",
		}
		if *json {
			// Use direct output
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		os.Exit(1)
	}

	// Check for shell alias
	aliasCmd := exec.Command("alias", command)
	aliasOutput, err := aliasCmd.Output()
	if err == nil && len(aliasOutput) > 0 {
		aliasStr := strings.TrimSpace(string(aliasOutput))
		aliasFile := findAliasFile(command)
		result := CheckAliasResult{
			Command: command,
			Type:    "alias",
			Value:   aliasStr,
			File:    aliasFile,
			Message: "Shell alias found",
		}
		if *json {
			// Use direct output
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		return
	}

	// Check for shell function
	funcCmd := exec.Command("type", command)
	funcOutput, err := funcCmd.Output()
	if err == nil && strings.Contains(string(funcOutput), "function") {
		funcFile := findFunctionFile(command)
		result := CheckAliasResult{
			Command: command,
			Type:    "function",
			Value:   "function " + command + "() { ... }",
			File:    funcFile,
			Message: "Shell function found",
		}
		if *json {
			// Use direct output
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		return
	}

	// Check for git alias
	gitCmd := exec.Command("git", "config", "--get", "alias."+command)
	gitOutput, err := gitCmd.Output()
	if err == nil && len(gitOutput) > 0 {
		gitAlias := strings.TrimSpace(string(gitOutput))
		result := CheckAliasResult{
			Command: command,
			Type:    "git_alias",
			Value:   "git alias " + command + "='" + gitAlias + "'",
			File:    "~/.gitconfig",
			Message: "Git alias found",
		}
		if *json {
			// Use direct output
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		return
	}

	// It's a binary/executable
	binaryPath := strings.TrimSpace(path)
	result := CheckAliasResult{
		Command: command,
		Type:    "binary",
		Value:   binaryPath,
		Message: "Binary executable found",
	}
	if *json {
		// Use direct output
		io.DirectOutput(result, *clip, *file, *json)
	} else {
		outputDefault(result)
	}
}

func findAliasFile(command string) string {
	// Common shell config files
	configFiles := []string{
		"~/.zshrc",
		"~/.zprofile",
		"~/.zshenv",
		"~/.zsh_aliases",
		"~/.aliases",
		"~/.bashrc",
		"~/.bash_aliases",
		"~/.profile",
		"~/.bash_profile",
	}

	// Add prompts-manager config files
	config, err := config.LoadConfig()
	if err == nil {
		managerDir := filepath.Dir(config.Prompts.BaseDir) // Go up one level from prompts dir
		configFiles = append(configFiles, filepath.Join(managerDir, "config/shell/*.zsh"))
		configFiles = append(configFiles, filepath.Join(managerDir, "config/shell/*.sh"))
	}

	// Search for alias definition
	for _, file := range configFiles {
		expandedFile := expandPath(file)
		if expandedFile != "" {
			cmd := exec.Command("grep", "-l", "alias "+command+"=", expandedFile)
			if err := cmd.Run(); err == nil {
				return file
			}
		}
	}

	return "unknown"
}

func findFunctionFile(command string) string {
	// Common shell config files
	configFiles := []string{
		"~/.zshrc",
		"~/.zprofile",
		"~/.zshenv",
		"~/.zsh_aliases",
		"~/.aliases",
		"~/.bashrc",
		"~/.bash_aliases",
		"~/.profile",
		"~/.bash_profile",
	}

	// Add prompts-manager config files
	config, err := config.LoadConfig()
	if err == nil {
		managerDir := filepath.Dir(config.Prompts.BaseDir) // Go up one level from prompts dir
		configFiles = append(configFiles, filepath.Join(managerDir, "config/shell/*.zsh"))
		configFiles = append(configFiles, filepath.Join(managerDir, "config/shell/*.sh"))
	}

	// Search for function definition
	for _, file := range configFiles {
		expandedFile := expandPath(file)
		if expandedFile != "" {
			cmd := exec.Command("grep", "-l", "^"+command+"()\\|^function "+command, expandedFile)
			if err := cmd.Run(); err == nil {
				return file
			}
		}
	}

	return "unknown"
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return strings.Replace(path, "~", home, 1)
	}
	return path
}

func outputDefault(result CheckAliasResult) {
	if result.Error != "" {
		fmt.Printf("❌ %s: %s\n", result.Command, result.Error)
		return
	}

	switch result.Type {
	case "alias":
		fmt.Printf("🔗 %s is a shell alias\n", result.Command)
		fmt.Printf("   Value: %s\n", result.Value)
		if result.File != "" {
			fmt.Printf("   File: %s\n", result.File)
		}
	case "function":
		fmt.Printf("⚙️  %s is a shell function\n", result.Command)
		fmt.Printf("   Value: %s\n", result.Value)
		if result.File != "" {
			fmt.Printf("   File: %s\n", result.File)
		}
	case "git_alias":
		fmt.Printf("📝 %s is a git alias\n", result.Command)
		fmt.Printf("   Value: %s\n", result.Value)
		if result.File != "" {
			fmt.Printf("   File: %s\n", result.File)
		}
	case "binary":
		fmt.Printf("🔧 %s is a binary executable\n", result.Command)
		fmt.Printf("   Path: %s\n", result.Value)
	case "not_found":
		fmt.Printf("❌ %s not found\n", result.Command)
	}
}
