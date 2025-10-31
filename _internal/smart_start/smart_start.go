package smart_start

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SmartStartInternalResult represents the result of a smart start operation
type SmartStartInternalResult struct {
	Action      string                 `json:"action"`
	RepoName    string                 `json:"repo_name"`
	ProjectType string                 `json:"project_type"`
	Project     string                 `json:"project,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	Success     bool                   `json:"success"`
	Error       string                 `json:"error,omitempty"`
	Message     string                 `json:"message"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// DetectAndConfigureProject detects project type and configures it
func DetectAndConfigureProject(repoName, repoRoot, directProject string) SmartStartInternalResult {
	// Simple: just execute different commands based on repo name
	switch repoName {
	case "rasch-stack":
		return startRaschStack(repoRoot)
	case "orbit":
		return startOrbit(repoRoot, directProject)
	default:
		// Check if it's a service (ends with -service)
		if strings.HasSuffix(repoName, "-service") {
			return startService(repoName, repoRoot)
		}
		// Check if it's a Go service
		if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
			return startGoService(repoName, repoRoot)
		}
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    repoName,
			ProjectType: "unknown",
			Success:     false,
			Error:       "Unknown project type",
			Message:     fmt.Sprintf("Unknown project type: %s", repoName),
		}
	}
}




func startRaschStack(repoRoot string) SmartStartInternalResult {
	// Check if ignition.js exists
	ignitionScript := filepath.Join(repoRoot, "scripts", "ignition.js")
	if _, err := os.Stat(ignitionScript); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    "rasch-stack",
			ProjectType: "rasch-stack",
			Success:     false,
			Error:       "ignition.js not found",
			Message:     "scripts/ignition.js not found in repository",
		}
	}

	// Start the rasch-stack project exactly like the original shell function
	// WDS_SOCKET_PORT=3333 NODE_OPTIONS=--openssl-legacy-provider node scripts/ignition.js --dev --nginx
	cmd := exec.Command("node", "scripts/ignition.js", "--dev", "--nginx")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), 
		"WDS_SOCKET_PORT=3333",
		"NODE_OPTIONS=--openssl-legacy-provider",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process in background
	if err := cmd.Start(); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    "rasch-stack",
			ProjectType: "rasch-stack",
			Success:     false,
			Error:       "Failed to start rasch-stack",
			Message:     fmt.Sprintf("Failed to start node scripts/ignition.js: %v", err),
		}
	}

	return SmartStartInternalResult{
		Action:      "_smart_start",
		RepoName:    "rasch-stack",
		ProjectType: "rasch-stack",
		Success:     true,
		Message:     "Rasch-stack project started (node scripts/ignition.js --dev --nginx)",
	}
}

func startOrbit(repoRoot string, directProject string) SmartStartInternalResult {
	// Check if package.json exists
	packageJson := filepath.Join(repoRoot, "package.json")
	if _, err := os.Stat(packageJson); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    "orbit",
			ProjectType: "orbit",
			Success:     false,
			Error:       "package.json not found",
			Message:     "package.json file not found in repository",
		}
	}

	// If no direct project provided, use a default
	if directProject == "" {
		directProject = "@dtc/orbit.rms.rocks" // Default project
	}

	// Start the orbit project exactly like the original shell function
	// ORBIT_SECRET_KEYS='{"basicAuthSAPService":"dummy"}' NODE_TLS_REJECT_UNAUTHORIZED=0 pnpm nx run "${sel}:local:dev"
	cmd := exec.Command("pnpm", "nx", "run", directProject+":local:dev")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(),
		"ORBIT_SECRET_KEYS={\"basicAuthSAPService\":\"dummy\"}",
		"NODE_TLS_REJECT_UNAUTHORIZED=0",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process in background
	if err := cmd.Start(); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    "orbit",
			ProjectType: "orbit",
			Success:     false,
			Error:       "Failed to start orbit project",
			Message:     fmt.Sprintf("Failed to start pnpm nx run: %v", err),
		}
	}

	return SmartStartInternalResult{
		Action:      "_smart_start",
		RepoName:    "orbit",
		ProjectType: "orbit",
		Success:     true,
		Message:     fmt.Sprintf("Orbit project started (pnpm nx run %s:local:dev)", directProject),
	}
}

func startGoService(repoName, repoRoot string) SmartStartInternalResult {
	// Check if go.mod and main.go exist
	goMod := filepath.Join(repoRoot, "go.mod")
	mainGo := filepath.Join(repoRoot, "main.go")
	
	if _, err := os.Stat(goMod); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    repoName,
			ProjectType: "service",
			Success:     false,
			Error:       "go.mod not found",
			Message:     "go.mod file not found in repository",
		}
	}
	
	if _, err := os.Stat(mainGo); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    repoName,
			ProjectType: "service",
			Success:     false,
			Error:       "main.go not found",
			Message:     "main.go file not found in repository",
		}
	}

	// Start the Go service exactly like the original shell function
	// go run main.go
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = repoRoot
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process in background
	if err := cmd.Start(); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    repoName,
			ProjectType: "service",
			Success:     false,
			Error:       "Failed to start Go service",
			Message:     fmt.Sprintf("Failed to start go run main.go: %v", err),
		}
	}

	return SmartStartInternalResult{
		Action:      "_smart_start",
		RepoName:    repoName,
		ProjectType: "service",
		Success:     true,
		Message:     "Go service started (go run main.go)",
	}
}

func startService(repoName, repoRoot string) SmartStartInternalResult {
	// Check if package.json exists (for yarn local:env)
	packageJson := filepath.Join(repoRoot, "package.json")
	if _, err := os.Stat(packageJson); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    repoName,
			ProjectType: "service",
			Success:     false,
			Error:       "package.json not found",
			Message:     "package.json file not found in repository",
		}
	}

	// Start the service project exactly like the original shell function
	// yarn local:dev
	cmd := exec.Command("yarn", "local:dev")
	cmd.Dir = repoRoot
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process in background
	if err := cmd.Start(); err != nil {
		return SmartStartInternalResult{
			Action:      "_smart_start",
			RepoName:    repoName,
			ProjectType: "service",
			Success:     false,
			Error:       "Failed to start service",
			Message:     fmt.Sprintf("Failed to start yarn local:dev: %v", err),
		}
	}

	return SmartStartInternalResult{
		Action:      "_smart_start",
		RepoName:    repoName,
		ProjectType: "service",
		Success:     true,
		Message:     "Service started (yarn local:dev)",
	}
}

func findPublications(srcDir string) ([]string, error) {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	var publications []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "common" && entry.Name() != "shared" {
			publications = append(publications, entry.Name())
		}
	}

	return publications, nil
}


