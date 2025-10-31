package main

// DESCRIPTION: apply zshrc

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"cli-go/_internal/ai"
)

func main() {

	flag.Parse()

	// Source the zshrc file using zsh (not sh)
	cmd := exec.Command("zsh", "-c", "source ~/.zshrc")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	ai.ExitIf(cmd.Run(), "failed to source zshrc")

	fmt.Fprintf(os.Stderr, "Zshrc configuration applied successfully\n")
}
