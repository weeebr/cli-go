package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PerplexityClient handles Perplexity API interactions
type PerplexityClient struct {
	apiKey string
	client *http.Client
}

// PerplexityRequest represents the API request structure
type PerplexityRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PerplexityResponse represents the API response structure
type PerplexityResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// Choice represents a response choice
type Choice struct {
	Message Message `json:"message"`
}

// Error represents an API error
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewPerplexityClient creates a new Perplexity client
func NewPerplexityClient(apiKey string) *PerplexityClient {
	return &PerplexityClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search performs a web search using Perplexity API
func (c *PerplexityClient) Search(query string) (string, error) {
	reqBody := PerplexityRequest{
		Model: "sonar-pro",
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.2,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://api.perplexity.ai/chat/completions",
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
		return "", c.handleAPIError(resp.StatusCode, body, query)
	}

	var response PerplexityResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response content received")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("empty response content")
	}

	return content, nil
}

// handleAPIError handles different API error scenarios with mock responses
func (c *PerplexityClient) handleAPIError(statusCode int, body []byte, query string) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return &MockAPIError{
			StatusCode: statusCode,
			Query:      query,
			Message:    "API key is invalid or expired",
			MockResponse: fmt.Sprintf(`Mock search result for: "%s"

This is a mock response because the API key is invalid or expired.
To get real results, update your API key:
  Run: setup

The system is working correctly - it just needs a valid API key to make real API calls.`, query),
		}
	case http.StatusTooManyRequests:
		return &MockAPIError{
			StatusCode: statusCode,
			Query:      query,
			Message:    "rate limit exceeded",
			MockResponse: fmt.Sprintf(`Mock search result for: "%s"

This is a mock response because the API rate limit was exceeded.
Please wait before making another request.

The system is working correctly - it just needs to wait before making more requests.`, query),
		}
	case http.StatusBadRequest:
		var response PerplexityResponse
		if err := json.Unmarshal(body, &response); err == nil && response.Error != nil {
			return fmt.Errorf("bad request: %s", response.Error.Message)
		}
		return fmt.Errorf("bad request - check your query format")
	case http.StatusInternalServerError:
		return &MockAPIError{
			StatusCode: statusCode,
			Query:      query,
			Message:    "server error",
			MockResponse: fmt.Sprintf(`Mock search result for: "%s"

This is a mock response because the Perplexity API server is experiencing issues.
Please try again later.

The system is working correctly - it just needs the API server to be available.`, query),
		}
	default:
		return &MockAPIError{
			StatusCode: statusCode,
			Query:      query,
			Message:    fmt.Sprintf("HTTP error %d", statusCode),
			MockResponse: fmt.Sprintf(`Mock search result for: "%s"

This is a mock response because the API request failed with status %d.
Please check your connection and try again.

The system is working correctly - it just needs a successful API response.`, query, statusCode),
		}
	}
}

// MockAPIError represents an API error that should return a mock response
type MockAPIError struct {
	StatusCode   int
	Query        string
	Message      string
	MockResponse string
}

func (e *MockAPIError) Error() string {
	return e.Message
}

// GetMockResponse returns the mock response content
func (e *MockAPIError) GetMockResponse() string {
	return e.MockResponse
}
