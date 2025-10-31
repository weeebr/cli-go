package main

// DESCRIPTION: Add individual API keys to credential store
// SECURITY: This tool is restricted to interactive use only
// SECURITY: No programmatic access allowed - requires user confirmation

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
	if !requireUserConfirmation("Adding API keys to credential store") {
		fmt.Println("❌ Operation cancelled by user")
		os.Exit(1)
	}

	fmt.Println("🔐 Add API Key to Credential Store")
	fmt.Println()

	store, err := config.NewEncryptedStore()
	ai.ExitIf(err, "failed to initialize credential store")

	// Load existing credentials
	var creds *config.Credentials
	if store.Exists() {
		creds, err = store.LoadCredentials()
		ai.ExitIf(err, "failed to load existing credentials")
	} else {
		// Create new credentials struct if none exist
		creds = &config.Credentials{}
	}

	// Display available key types
	fmt.Println("Available API key types:")
	fmt.Println("1. OpenAI")
	fmt.Println("2. Anthropic")
	fmt.Println("3. Google")
	fmt.Println("4. Groq")
	fmt.Println("5. Perplexity")
	fmt.Println("6. Figma")
	fmt.Println("7. Jira API Token")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select key type (1-7): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var keyName string
	var promptText string

	switch choice {
	case "1":
		keyName = "OpenAI"
		promptText = "OpenAI API Key: "
	case "2":
		keyName = "Anthropic"
		promptText = "Anthropic API Key: "
	case "3":
		keyName = "Google"
		promptText = "Google API Key: "
	case "4":
		keyName = "Groq"
		promptText = "Groq API Key: "
	case "5":
		keyName = "Perplexity"
		promptText = "Perplexity API Key: "
	case "6":
		keyName = "Figma"
		promptText = "Figma API Token: "
	case "7":
		keyName = "Jira API Token"
		promptText = "Jira API Token: "
	default:
		fmt.Fprintf(os.Stderr, "❌ Invalid choice: %s\n", choice)
		os.Exit(1)
	}

	// Get the key value (visible input for verification)
	fmt.Print(promptText)
	keyValueStr, _ := reader.ReadString('\n')
	keyValueStr = strings.TrimSpace(keyValueStr)

	if keyValueStr == "" {
		fmt.Fprintf(os.Stderr, "❌ %s key cannot be empty\n", keyName)
		os.Exit(1)
	}

	// Update the appropriate credential field
	switch choice {
	case "1":
		creds.OpenAI = keyValueStr
	case "2":
		creds.Anthropic = keyValueStr
	case "3":
		creds.Google = keyValueStr
	case "4":
		creds.Groq = keyValueStr
	case "5":
		creds.Perplexity = keyValueStr
	case "6":
		creds.Figma = keyValueStr
	case "7":
		creds.Jira = keyValueStr
	}

	// Save the updated credentials
	ai.ExitIf(store.SaveCredentials(creds), "failed to save credentials")

	fmt.Printf("\n✅ Successfully added %s key\n", keyName)
	fmt.Println("🔒 Credentials saved to ~/.cli-go/credentials.enc")
	fmt.Println("⚠️  To update: re-run this tool or use 'setup' for full reconfiguration")
}
