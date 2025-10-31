package sys

import (
	"os"
)

// PackageManager represents a package manager
type PackageManager string

const (
	ManagerYarn PackageManager = "yarn"
	ManagerPnpm PackageManager = "pnpm"
	ManagerNpm  PackageManager = "npm"
	ManagerNone PackageManager = "none"
)

// DetectPackageManager detects which package manager to use
func DetectPackageManager() PackageManager {
	// Check for lock files in order of preference
	if _, err := os.Stat("yarn.lock"); err == nil {
		return ManagerYarn
	}
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return ManagerPnpm
	}
	if _, err := os.Stat("package-lock.json"); err == nil {
		return ManagerNpm
	}

	// Check for package.json to determine if it's a Node.js project
	if _, err := os.Stat("package.json"); err == nil {
		return ManagerNpm // Default to npm if package.json exists
	}

	return ManagerNone
}

// InstallResult holds the result of an installation
type InstallResult struct {
	Installed bool           `json:"installed"`
	Manager   PackageManager `json:"manager"`
	Error     string         `json:"error,omitempty"`
}

// InstallPackages installs packages using the detected manager
func InstallPackages() *InstallResult {
	manager := DetectPackageManager()

	if manager == ManagerNone {
		return &InstallResult{
			Installed: false,
			Manager:   manager,
			Error:     "no package.json found",
		}
	}

	var cmd string
	var args []string

	switch manager {
	case ManagerYarn:
		cmd = "yarn"
		args = []string{"install"}
	case ManagerPnpm:
		cmd = "pnpm"
		args = []string{"install"}
	case ManagerNpm:
		cmd = "npm"
		args = []string{"install", "--legacy-peer-deps"}
	}

	result := RunCommand(cmd, args...)
	if result.ExitCode != 0 {
		return &InstallResult{
			Installed: false,
			Manager:   manager,
			Error:     result.Stderr,
		}
	}

	return &InstallResult{
		Installed: true,
		Manager:   manager,
	}
}

// ReinstallPackages removes node_modules and reinstalls
func ReinstallPackages() *InstallResult {
	// Remove node_modules if it exists
	if _, err := os.Stat("node_modules"); err == nil {
		os.RemoveAll("node_modules")
	}

	// Install packages
	return InstallPackages()
}
