package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"cli-go/_internal/config"
	"io"
	"net/http"
	"time"
)

// GeminiClient handles Gemini HTTP API integration
type GeminiClient struct {
	apiKey string
	client *http.Client
}

// GeminiRequest represents the Google Gemini API request structure
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content in the request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of the content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents the Google Gemini API response structure
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
	Error      *GeminiError      `json:"error,omitempty"`
}

// GeminiCandidate represents a response candidate
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// GeminiError represents an API error
type GeminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(apiKey string) *GeminiClient {
	// Load config to get timeout setting
	cfg, err := config.LoadConfig()
	timeout := 60 * time.Second // Default fallback
	if err == nil {
		timeout = time.Duration(cfg.AI.Timeouts.Default) * time.Second
	}

	return &GeminiClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// NewGeminiClientFromStore creates a new Gemini client using encrypted store
func NewGeminiClientFromStore() (*GeminiClient, error) {
	apiKey, err := config.GetKey("google")
	if err != nil {
		return nil, fmt.Errorf("failed to get Google API key: %v", err)
	}

	return NewGeminiClient(apiKey), nil
}

// SendMessage sends a message to Gemini via HTTP API
func (g *GeminiClient) SendMessage(message string) (string, error) {
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: message,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Load config to get model name
	cfg, err := config.LoadConfig()
	model := "gemini-2.5-flash" // Default fallback
	if err == nil {
		model = cfg.AI.Models.Google
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, g.apiKey)

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
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

	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s", response.Error.Message)
	}

	if len(response.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates received")
	}

	if len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in response")
	}

	content := response.Candidates[0].Content.Parts[0].Text
	if content == "" {
		return "", fmt.Errorf("empty response content")
	}

	return content, nil
}

// GetModel returns the model name
func (g *GeminiClient) GetModel() string {
	// Load config to get model name
	cfg, err := config.LoadConfig()
	if err == nil {
		return cfg.AI.Models.Google
	}
	return "gemini-2.5-flash" // Default fallback
}
