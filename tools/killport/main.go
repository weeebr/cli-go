package main

import (
	"flag"
	"fmt"
	"cli-go/_internal/flags"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var (
		explain = flag.Bool("explain", false, "Show detailed explanation")
	)
	flags.ReorderAndParse()

	if *explain {
		fmt.Fprintf(os.Stderr, "killport: Kill processes running on a specific port\n")
		fmt.Fprintf(os.Stderr, "Usage: killport <port>\n")
		fmt.Fprintf(os.Stderr, "Example: killport 3000\n")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: port required\n")
		fmt.Fprintf(os.Stderr, "Usage: killport <port>\n")
		os.Exit(1)
	}

	port, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid port number\n")
		fmt.Fprintf(os.Stderr, "Port must be a valid integer\n")
		os.Exit(1)
	}

	// Find processes using the port
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf("tcp:%d", port))
	output_bytes, err := cmd.Output()
	if err != nil {
		// lsof returns exit code 1 when no processes are found, which is normal
		fmt.Fprintf(os.Stderr, "No processes found on port %d\n", port)
		os.Exit(0)
	}

	// Parse PIDs
	pidLines := strings.TrimSpace(string(output_bytes))
	if pidLines == "" {
		fmt.Fprintf(os.Stderr, "No processes found on port %d\n", port)
		os.Exit(0)
	}

	pids := strings.Split(pidLines, "\n")
	var killedPids []int

	// Kill each process
	for _, pidStr := range pids {
		if pidStr == "" {
			continue
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Try graceful termination first (SIGTERM)
		killCmd := exec.Command("kill", "-TERM", pidStr)
		if err := killCmd.Run(); err != nil {
			// If graceful termination fails, try SIGKILL as fallback
			killCmd = exec.Command("kill", "-9", pidStr)
			if err := killCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to kill process %d: %v\n", pid, err)
				continue
			}
		}
		killedPids = append(killedPids, pid)
		fmt.Fprintf(os.Stderr, "Killed process %d\n", pid)
	}

	fmt.Fprintf(os.Stderr, "Killed %d processes on port %d\n", len(killedPids), port)
}
