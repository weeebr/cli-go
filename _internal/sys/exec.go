package sys

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecResult holds the result of a command execution
type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// RunCommand executes a command and returns the result
func RunCommand(name string, args ...string) *ExecResult {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return &ExecResult{
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
		ExitCode: exitCode,
		Error:    err,
	}
}

// RunCommandInDir executes a command in a specific directory
func RunCommandInDir(dir, name string, args ...string) *ExecResult {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return &ExecResult{
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
		ExitCode: exitCode,
		Error:    err,
	}
}

// CheckDependency checks if a single command exists
func CheckDependency(cmd string) error {
	_, err := exec.LookPath(cmd)
	return err
}

// CheckDependencies checks multiple commands
func CheckDependencies(deps []string) error {
	for _, dep := range deps {
		if err := CheckDependency(dep); err != nil {
			return fmt.Errorf("required dependency not found: %s", dep)
		}
	}
	return nil
}

// GetCurrentWorkingDirectory gets CWD
func GetCurrentWorkingDirectory() (string, error) {
	return os.Getwd()
}

// IsValidPath checks if path exists
func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
