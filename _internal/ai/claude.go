package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cli-go/_internal/config"
)

// ClaudeModel represents Claude model types
type ClaudeModel string

const (
	ClaudeSonnet ClaudeModel = "sonnet"
	ClaudeHaiku  ClaudeModel = "haiku"
)

// ClaudeClient handles Claude HTTP API integration
type ClaudeClient struct {
	model  ClaudeModel
	apiKey string
	client *http.Client
}

// ClaudeRequest represents the Anthropic API request structure
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeMessage represents a message in the conversation
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the Anthropic API response structure
type ClaudeResponse struct {
	ID      string               `json:"id"`
	Type    string               `json:"type"`
	Role    string               `json:"role"`
	Content []ClaudeContentBlock `json:"content"`
	Error   *ClaudeError         `json:"error,omitempty"`
}

// ClaudeContentBlock represents a content block in the response
type ClaudeContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeError represents an API error
type ClaudeError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// NewClaudeClient creates a new Claude client
func NewClaudeClient(model ClaudeModel, apiKey string) *ClaudeClient {
	// Load config to get timeout setting
	cfg, err := config.LoadConfig()
	timeout := 60 * time.Second // Default fallback
	if err == nil {
		timeout = time.Duration(cfg.AI.Timeouts.Default) * time.Second
	}

	return &ClaudeClient{
		model:  model,
		apiKey: apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// SendMessage sends a message to Claude via HTTP API
func (c *ClaudeClient) SendMessage(message string) (string, error) {
	// Load config to get model names
	cfg, err := config.LoadConfig()
	modelName := "claude-sonnet-4-5-20250929" // Default fallback
	if err == nil {
		if c.model == ClaudeHaiku {
			modelName = "claude-haiku-4-5-20251001" // Default haiku fallback
		} else {
			modelName = cfg.AI.Models.Anthropic
		}
	} else if c.model == ClaudeHaiku {
		modelName = "claude-haiku-4-5-20251001"
	}

	reqBody := ClaudeRequest{
		Model:     modelName,
		MaxTokens: 4096,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://api.anthropic.com/v1/messages",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

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

	var response ClaudeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("Claude API error: %s", response.Error.Message)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no response content received")
	}

	content := response.Content[0].Text
	if content == "" {
		return "", fmt.Errorf("empty response content")
	}

	return content, nil
}

// GetModel returns the current model
func (c *ClaudeClient) GetModel() string {
	// Load config to get model names
	cfg, err := config.LoadConfig()
	if err == nil {
		if c.model == ClaudeHaiku {
			return "claude-haiku-4-5-20251001" // Default haiku fallback
		}
		return cfg.AI.Models.Anthropic
	}

	// Fallback to hardcoded values
	if c.model == ClaudeHaiku {
		return "claude-haiku-4-5-20251001"
	}
	return "claude-sonnet-4-5-20250929"
}
