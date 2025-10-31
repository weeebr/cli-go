package main

// DESCRIPTION: measure $1's time (5x)

import (
	"context"
	"flag"
	"fmt"
	"cli-go/_internal/flags"
	"cli-go/_internal/io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PerfResult struct {
	Command   string    `json:"command"`
	Runs      int       `json:"runs"`
	TotalTime float64   `json:"total_time_seconds"`
	AvgTime   float64   `json:"avg_time_seconds"`
	Times     []float64 `json:"individual_times,omitempty"`
	Timeout   int       `json:"timeout_seconds,omitempty"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message"`
}

func main() {

	var (
		clip    = flag.Bool("clip", false, "Copy to clipboard")
		file    = flag.String("file", "", "Write to file")
		json    = flag.Bool("json", false, "Output as JSON (for piping)")
		timeout = flag.Int("timeout", 30, "Timeout in seconds for each command execution")
	)
	flags.ReorderAndParse()

	args := flag.Args()
	if len(args) == 0 {
		result := PerfResult{
			Error:   "command required",
			Message: "Usage: perf 'command | with | pipes'",
		}
		if *json {
			// Use centralized output routing
			io.DirectOutput(result, *clip, *file, *json)
		} else {
			outputDefault(result)
		}
		os.Exit(1)
	}

	command := strings.Join(args, " ")
	const runs = 5

	fmt.Fprintf(os.Stderr, "Testing: %s\n", command)
	fmt.Fprintf(os.Stderr, "Running %d times...\n", runs)

	var times []float64
	totalStart := time.Now()

	for i := 0; i < runs; i++ {
		start := time.Now()

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
		defer cancel()

		// Execute the command with timeout
		cmd := exec.CommandContext(ctx, "sh", "-c", command)
		cmd.Stdout = os.Stderr // Redirect output to stderr to avoid mixing with JSON
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		elapsed := time.Since(start).Seconds()

		if err != nil {
			var errorMsg string
			if ctx.Err() == context.DeadlineExceeded {
				errorMsg = fmt.Sprintf("command timed out after %ds: %v", *timeout, err)
			} else {
				errorMsg = fmt.Sprintf("command failed: %v", err)
			}

			result := PerfResult{
				Command: command,
				Timeout: *timeout,
				Error:   errorMsg,
				Message: "Command execution failed",
			}
			if *json {
				// Use centralized output routing
				io.DirectOutput(result, *clip, *file, *json)
			} else {
				outputDefault(result)
			}
			os.Exit(1)
		}

		times = append(times, elapsed)
	}

	totalTime := time.Since(totalStart).Seconds()
	avgTime := totalTime / float64(runs)

	result := PerfResult{
		Command:   command,
		Runs:      runs,
		TotalTime: totalTime,
		AvgTime:   avgTime,
		Times:     times,
		Timeout:   *timeout,
		Message:   fmt.Sprintf("Total: %.3fs | Average: %.3fs", totalTime, avgTime),
	}

	if *json {
		// Use centralized output routing
		io.DirectOutput(result, *clip, *file, *json)
	} else {
		outputDefault(result)
	}
}

func formatMarkdownResult(result PerfResult) string {
	if result.Error != "" {
		return fmt.Sprintf("❌ %s", result.Error)
	} else {
		return fmt.Sprintf("✅ %s", result.Message)
	}
}

func outputDefault(result PerfResult) {
	if result.Error != "" {
		fmt.Printf("❌ %s: %s\n", result.Command, result.Error)
		return
	}

	fmt.Printf("⏱️  Performance results for: %s\n", result.Command)
	fmt.Printf("   Runs: %d\n", result.Runs)
	fmt.Printf("   Timeout: %ds\n", result.Timeout)
	fmt.Printf("   Total time: %.3fs\n", result.TotalTime)
	fmt.Printf("   Average time: %.3fs\n", result.AvgTime)

	if len(result.Times) > 0 {
		fmt.Printf("   Individual times: ")
		for i, time := range result.Times {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%.3fs", time)
		}
		fmt.Printf("\n")
	}
}
