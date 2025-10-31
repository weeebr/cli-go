package main

// DESCRIPTION: CLI tools config

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"cli-go/_internal/ai"
	"cli-go/_internal/config"
)

func main() {

	var (
		editor = flag.String("editor", "textedit", "Editor to use (textedit, cursor, code, vim)")
	)
	flag.Parse()

	// Get config file path relative to cli-go tool location
	configFile, err := config.GetConfigPath()
	ai.ExitIf(err, "failed to get config path")

	// Check if config file exists
	_, err = os.Stat(configFile)
	ai.ExitIf(err, fmt.Sprintf("config.yml not found in cli-go directory: %s", configFile))

	// Open config file in editor
	var editCmd *exec.Cmd
	switch *editor {
	case "textedit":
		editCmd = exec.Command("open", "-a", "TextEdit", configFile)
	case "cursor":
		editCmd = exec.Command("cursor", configFile)
	case "code":
		editCmd = exec.Command("code", configFile)
	case "vim":
		editCmd = exec.Command("vim", configFile)
	default:
		editCmd = exec.Command(*editor, configFile)
	}

	editCmd.Stdout = os.Stderr
	editCmd.Stderr = os.Stderr

	ai.ExitIf(editCmd.Run(), fmt.Sprintf("failed to open config in %s", *editor))

	fmt.Fprintf(os.Stderr, "Opened config file in %s\n", *editor)
}
