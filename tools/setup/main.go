package main

// DESCRIPTION: Go CLI setup - goes through all 10 key types

import (
	"bufio"
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
	"os"
	"strings"

	"golang.org/x/term"
)

// SECURITY GATE: Interactive confirmation required
func requireUserConfirmation(action string) bool {
	fmt.Printf("⚠️  SECURITY: %s requires user confirmation\n", action)
	fmt.Print("Press Enter to continue: ")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	return true
}

// SECURITY GATE: Prevent non-interactive usage
func checkInteractiveMode() {
	// Check if running in interactive terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(os.Stderr, "❌ SECURITY: This tool requires interactive terminal")
		fmt.Fprintln(os.Stderr, "❌ Non-interactive usage is blocked for security")
		os.Exit(1)
	}

	// Additional check: ensure we have a controlling terminal
	if os.Getenv("TERM") == "" {
		fmt.Fprintln(os.Stderr, "❌ SECURITY: No terminal environment detected")
		os.Exit(1)
	}
}

func main() {
	// SECURITY: Enforce interactive mode only
	checkInteractiveMode()

	// SECURITY: Require user confirmation for credential operations
	if !requireUserConfirmation("Setting up all API keys") {
		fmt.Println("❌ Operation cancelled by user")
		os.Exit(1)
	}

	fmt.Println("🔐 Go CLI Complete Credential Setup")
	fmt.Println("Going through all 10 key types...")
	fmt.Println()

	store, err := config.NewEncryptedStore()
	ai.ExitIf(err, "failed to initialize credential store")

	// Load existing credentials or create new
	var creds *config.Credentials
	if store.Exists() {
		creds, err = store.LoadCredentials()
		ai.ExitIf(err, "failed to load existing credentials")
	} else {
		creds = &config.Credentials{}
	}

	reader := bufio.NewReader(os.Stdin)

	// Go through all 10 key types
	fmt.Println("Enter API keys (press Enter to skip):")

	creds.OpenAI = promptKey(reader, "OpenAI: ")
	creds.Anthropic = promptKey(reader, "Anthropic: ")
	creds.Google = promptKey(reader, "Google (optional): ")
	creds.Groq = promptKey(reader, "Groq (optional): ")
	creds.Perplexity = promptKey(reader, "Perplexity (optional): ")
	creds.Figma = promptKey(reader, "Figma (optional): ")

	creds.Jira = promptKey(reader, "Jira: ")

	// Check if at least one key is provided
	if creds.OpenAI == "" && creds.Anthropic == "" && creds.Google == "" && creds.Groq == "" && creds.Perplexity == "" && creds.Figma == "" {
		fmt.Fprintln(os.Stderr, "❌ At least one API key required")
		os.Exit(1)
	}

	// Save the updated credentials
	ai.ExitIf(store.SaveCredentials(creds), "failed to save credentials")

	fmt.Println("\n✅ All credentials saved to ~/.cli-go/credentials.enc")
	fmt.Println("💡 Use './bin/add_key' to update individual keys")
}

func promptKey(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	return strings.TrimSpace(value)
}
