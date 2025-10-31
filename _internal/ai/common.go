package ai

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// InputMode represents how input is provided to AI tools
type InputMode string

const (
	InputArgs        InputMode = "args"
	InputStdin       InputMode = "stdin"
	InputInteractive InputMode = "interactive"
)

// DetectInputMode determines how input is being provided
func DetectInputMode() InputMode {
	// Check if there are non-flag command line arguments
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if !strings.HasPrefix(arg, "-") {
			return InputArgs
		}
	}

	// Check if stdin is a terminal
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return InputStdin
	}

	return InputInteractive
}

// ReadStdin reads all input from stdin
func ReadStdin() (string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

// GetArgs returns command line arguments as a single string
func GetArgs() string {
	args := flag.Args()
	if len(args) == 0 {
		return ""
	}
	return strings.Join(args, " ")
}

// CheckError detects common AI CLI errors
func CheckError(output string) error {
	output = strings.ToLower(output)
	if strings.Contains(output, "error") ||
		strings.Contains(output, "failed") ||
		strings.Contains(output, "quota") {
		return fmt.Errorf("AI CLI error detected")
	}
	return nil
}

// LogError writes error message to stderr
func LogError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// LogInfo writes info message to stderr
func LogInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// ResponseInfo holds response metadata for AI tools
type ResponseInfo struct {
	Duration time.Duration
	Model    string
}

// FormatResponseInfo formats response time and model info with emoji
func FormatResponseInfo(info ResponseInfo) string {
	duration := info.Duration
	var durationStr string

	if duration < time.Millisecond {
		durationStr = fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1e6)
	} else if duration < time.Second {
		durationStr = fmt.Sprintf("%.2fs", duration.Seconds())
	} else {
		durationStr = fmt.Sprintf("%.1fs", duration.Seconds())
	}

	return fmt.Sprintf("🤖 %s, Model: %s", durationStr, info.Model)
}

// AIClient interface for AI clients that can track response time
type AIClient interface {
	SendMessage(message string) (string, error)
	GetModel() string
}

// SendMessageWithTiming wraps AI client calls with response time tracking
func SendMessageWithTiming(client AIClient, message string) (string, ResponseInfo, error) {
	start := time.Now()
	response, err := client.SendMessage(message)
	duration := time.Since(start)

	info := ResponseInfo{
		Duration: duration,
		Model:    client.GetModel(),
	}

	return response, info, err
}
