package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"cli-go/_internal/config"
)

// ChatGPTConfig holds configuration for ChatGPT integration
type ChatGPTConfig struct {
	Thread string
	Format string // "text" or "json"
	Model  string // GPT model to use
}

// ChatGPTClient handles ChatGPT HTTP API integration
type ChatGPTClient struct {
	config ChatGPTConfig
	apiKey string
	client *http.Client
}

// ChatGPTRequest represents the OpenAI API request structure
type ChatGPTRequest struct {
	Model       string           `json:"model"`
	Messages    []ChatGPTMessage `json:"messages"`
	Temperature float64          `json:"temperature,omitempty"`
}

// ChatGPTMessage represents a message in the conversation
type ChatGPTMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp,omitempty"`
}

// ChatGPTResponse represents the OpenAI API response structure
type ChatGPTResponse struct {
	Choices []ChatGPTChoice `json:"choices"`
	Error   *ChatGPTError   `json:"error,omitempty"`
}

// ChatGPTChoice represents a response choice
type ChatGPTChoice struct {
	Message ChatGPTMessage `json:"message"`
}

// ChatGPTError represents an API error
type ChatGPTError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// NewChatGPTClient creates a new ChatGPT client
func NewChatGPTClient(format string, apiKey string) *ChatGPTClient {
	// Load config to get model name
	cfg, err := config.LoadConfig()
	model := "gpt-4o" // Default fallback
	if err == nil {
		model = cfg.AI.Models.OpenAI
	}

	return &ChatGPTClient{
		config: ChatGPTConfig{
			Thread: getCurrentThread(),
			Format: format,
			Model:  model,
		},
		apiKey: apiKey,
		client: &http.Client{
			Timeout: time.Duration(cfg.AI.Timeouts.Default) * time.Second,
		},
	}
}

// NewChatGPTClientWithModel creates a new ChatGPT client with specific model
func NewChatGPTClientWithModel(format, model, apiKey string) *ChatGPTClient {
	return &ChatGPTClient{
		config: ChatGPTConfig{
			Thread: getCurrentThread(),
			Format: format,
			Model:  model,
		},
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// SendMessage sends a message to ChatGPT via HTTP API
func (c *ChatGPTClient) SendMessage(message string) (string, error) {
	// Load conversation history
	history, err := c.loadThreadHistory()
	if err != nil {
		return "", err
	}

	// Add user message to history
	userMessage := ChatGPTMessage{
		Role:      "user",
		Content:   message,
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
	}
	history = append(history, userMessage)

	// Prepare request
	reqBody := ChatGPTRequest{
		Model:    c.config.Model,
		Messages: history,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response ChatGPTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("ChatGPT API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices received")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("empty response content")
	}

	// Add assistant message to history
	assistantMessage := ChatGPTMessage{
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
	}
	history = append(history, assistantMessage)

	// Save updated history
	if err := c.saveThreadHistory(history); err != nil {
		LogError("Failed to save thread history: %v", err)
	}

	return content, nil
}

// SendMessageWithRoleFile sends a message using a role file
func (c *ChatGPTClient) SendMessageWithRoleFile(roleFile, message string) (string, error) {
	// Read role file
	roleData, err := os.ReadFile(roleFile)
	if err != nil {
		return "", fmt.Errorf("failed to read role file: %v", err)
	}

	// Load conversation history
	history, err := c.loadThreadHistory()
	if err != nil {
		return "", err
	}

	// Replace system message with role file content
	if len(history) > 0 && history[0].Role == "system" {
		history[0].Content = string(roleData)
	} else {
		// Prepend system message if it doesn't exist
		systemMessage := ChatGPTMessage{
			Role:      "system",
			Content:   string(roleData),
			Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
		}
		history = append([]ChatGPTMessage{systemMessage}, history...)
	}

	// Add user message to history
	userMessage := ChatGPTMessage{
		Role:      "user",
		Content:   message,
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
	}
	history = append(history, userMessage)

	// Prepare request
	reqBody := ChatGPTRequest{
		Model:    c.config.Model,
		Messages: history,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response ChatGPTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("ChatGPT API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices received")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("empty response content")
	}

	// Add assistant message to history
	assistantMessage := ChatGPTMessage{
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
	}
	history = append(history, assistantMessage)

	// Save updated history
	if err := c.saveThreadHistory(history); err != nil {
		LogError("Failed to save thread history: %v", err)
	}

	return content, nil
}

// GetThread returns the current thread name
func (c *ChatGPTClient) GetThread() string {
	return c.config.Thread
}

// GetModel returns the model name
func (c *ChatGPTClient) GetModel() string {
	return c.config.Model
}
