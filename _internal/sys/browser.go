package sys

import (
	"os/exec"
	"runtime"
)

// OpenInBrowser opens a URL in the default browser
// Supports Windows, macOS, and Linux
func OpenInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		// Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Run()
}
