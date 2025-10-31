package main

// DESCRIPTION: edit config (if no args, otherwise run)

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"cli-go/_internal/ai"
)

func main() {

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		// No arguments - open zshrc in TextEdit
		cmd := exec.Command("open", "-a", "TextEdit", "~/.zshrc")
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr

		ai.ExitIf(cmd.Run(), "failed to open zshrc")

		fmt.Fprintf(os.Stderr, "Opened ~/.zshrc in TextEdit\n")
		return
	}

	// Has arguments - run zsh with those arguments
	cmd := exec.Command("zsh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run zsh interactively
	ai.ExitIf(cmd.Run(), "zsh execution failed")
}
