package main

// DESCRIPTION: Check current credentials in encrypted store

import (
	"fmt"
	"cli-go/_internal/ai"
	"cli-go/_internal/config"
)

func main() {
	store, err := config.NewEncryptedStore()
	ai.ExitIf(err, "failed to initialize credential store")

	if !store.Exists() {
		fmt.Println("📭 No credentials file found")
		fmt.Println("💡 Run './bin/setup' or './bin/add_key' to add credentials")
		return
	}

	creds, err := store.LoadCredentials()
	ai.ExitIf(err, "failed to load credentials")

	fmt.Println("🔐 Current Credentials in Encrypted Store:")
	fmt.Println()

	// Check each credential type
	checkCred("OpenAI", creds.OpenAI)
	checkCred("Anthropic", creds.Anthropic)
	checkCred("Google", creds.Google)
	checkCred("Groq", creds.Groq)
	checkCred("Perplexity", creds.Perplexity)
	checkCred("Figma", creds.Figma)

	checkCred("Jira", creds.Jira)
}

func checkCred(name, value string) {
	if value != "" {
		// Show first 8 chars + ... for security
		masked := value[:8] + "..."
		fmt.Printf("✅ %s: %s\n", name, masked)
	} else {
		fmt.Printf("❌ %s: Not configured\n", name)
	}
}
