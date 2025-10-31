package flags

import (
	"flag"
	"os"
	"strings"
)

// ReorderAndParse separates flags from positional args and parses them correctly
// This allows flags to come before or after positional arguments
// Supports X = -X = --X format (single letter without dash gets auto-prefixed)
func ReorderAndParse() {
	var flags []string
	var args []string

	// Process all arguments except the program name
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-") {
			// Already a flag
			flags = append(flags, arg)
			// Check if next arg is a flag value (not starting with -)
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") && !isBoolFlag(arg) {
				i++
				flags = append(flags, os.Args[i])
			}
		} else if len(arg) == 1 && isBoolFlag("-"+arg) {
			// Single letter that matches a known bool flag - treat as flag
			flags = append(flags, "-"+arg)
		} else {
			// Positional argument
			args = append(args, arg)
		}
	}

	// Reconstruct os.Args: [program, flags..., args...]
	os.Args = append([]string{os.Args[0]}, append(flags, args...)...)
	flag.Parse()
}

// isBoolFlag checks if a flag is a boolean flag (doesn't take a value)
// Common boolean flags in the CLI tools
func isBoolFlag(flag string) bool {
	boolFlags := map[string]bool{
		"-h":        true,
		"--json":    true,
		"--compact": true,
		"--all":     true,
		"--main":    true,
		"--current": true,
		"-o":        true,
		"--o":       true,
		"--open":    true,
		"--verbose": true,
		"-v":        true,
		"--version": true,
		"--force":   true,
		"-f":        true,
		"--dry-run": true,
		"--quiet":   true,
		"-q":        true,
		"--yes":     true,
		"-y":        true,
		"--no":      true,
		"-n":        true,
	}

	// Remove leading dashes for lookup
	cleanFlag := strings.TrimLeft(flag, "-")
	return boolFlags[flag] || boolFlags["--"+cleanFlag] || boolFlags["-"+cleanFlag]
}
