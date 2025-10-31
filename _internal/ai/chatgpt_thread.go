package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"cli-go/_internal/config"
)

// getHistoryDir returns the history directory from config
func getHistoryDir() string {
	cfg, err := config.LoadConfig()
	if err == nil {
		// Use config cache directory
		baseDir := cfg.Cache.BaseDir
		if strings.HasPrefix(baseDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				baseDir = strings.Replace(baseDir, "~", homeDir, 1)
			}
		}
		return filepath.Join(baseDir, "chatgpt", "history")
	}
	// Fallback to default
	return filepath.Join(os.Getenv("HOME"), ".chatgpt-cli", "history")
}

// getCurrentThread gets or creates the current thread based on shell PGID
func getCurrentThread() string {
	// Get current process group ID
	pgid := syscall.Getpgrp()

	// Create thread name based on date and PGID
	baseDate := time.Now().Format("02-01-2006")
	thread := fmt.Sprintf("%s-%d", baseDate, pgid)

	// Ensure thread file exists
	createThreadFile(thread)

	return thread
}

// createThreadFile creates the initial thread file
func createThreadFile(thread string) {
	historyDir := getHistoryDir()
	os.MkdirAll(historyDir, 0755)

	threadFile := filepath.Join(historyDir, thread+".json")
	if _, err := os.Stat(threadFile); os.IsNotExist(err) {
		initialData := []ChatGPTMessage{
			{
				Role:      "system",
				Content:   "You are a helpful assistant.",
				Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
			},
		}

		data, _ := json.MarshalIndent(initialData, "", "  ")
		os.WriteFile(threadFile, data, 0644)
	}
}

// loadThreadHistory loads the conversation history from the thread file
func (c *ChatGPTClient) loadThreadHistory() ([]ChatGPTMessage, error) {
	historyDir := getHistoryDir()
	threadFile := filepath.Join(historyDir, c.config.Thread+".json")

	data, err := os.ReadFile(threadFile)
	if err != nil {
		// If thread file doesn't exist, create it
		createThreadFile(c.config.Thread)
		data, err = os.ReadFile(threadFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read thread file: %v", err)
		}
	}

	var messages []ChatGPTMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to parse thread history: %v", err)
	}

	return messages, nil
}

// saveThreadHistory saves the conversation history to the thread file
func (c *ChatGPTClient) saveThreadHistory(messages []ChatGPTMessage) error {
	historyDir := getHistoryDir()
	threadFile := filepath.Join(historyDir, c.config.Thread+".json")

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal thread history: %v", err)
	}

	if err := os.WriteFile(threadFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write thread file: %v", err)
	}

	return nil
}
