package sys

import (
	"os/exec"
	"runtime"
	"strings"
)

// CopyToClipboard copies content to the system clipboard
// Uses cross-platform approach with fallback to platform-specific commands
func CopyToClipboard(content string) error {
	// Try cross-platform library first (if available)
	// For now, use platform-specific commands for reliability

	switch runtime.GOOS {
	case "darwin":
		// macOS - use pbcopy
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(content)
		return cmd.Run()
	case "windows":
		// Windows - use clip command
		cmd := exec.Command("clip")
		cmd.Stdin = strings.NewReader(content)
		return cmd.Run()
	case "linux":
		// Linux - try xclip first, then xsel
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(content)
		if err := cmd.Run(); err != nil {
			// Fallback to xsel
			cmd = exec.Command("xsel", "--clipboard", "--input")
			cmd.Stdin = strings.NewReader(content)
			return cmd.Run()
		}
		return nil
	default:
		// Unknown platform - try pbcopy as fallback (common on macOS-like systems)
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(content)
		return cmd.Run()
	}
}
