package figma

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cli-go/_internal/config"
)

// Client represents a Figma API client
type Client struct {
	Token      string
	FileKey    string
	HTTPClient *http.Client
}

// NewClient creates a new Figma API client
func NewClient(token, fileKey string) *Client {
	// Load config to get timeout setting
	cfg, err := config.LoadConfig()
	timeout := 30 * time.Second // Default fallback
	if err == nil {
		timeout = time.Duration(cfg.Network.TimeoutSeconds) * time.Second
	}

	return &Client{
		Token:   token,
		FileKey: fileKey,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ValidateToken checks if token has correct format (figd_ prefix)
func (c *Client) ValidateToken() error {
	if c.Token == "" {
		return fmt.Errorf("FIGMA_API_TOKEN required")
	}
	if !strings.HasPrefix(c.Token, "figd_") {
		return fmt.Errorf("FIGMA_API_TOKEN must start with 'figd_'")
	}
	return nil
}

// FetchFile fetches the entire Figma file structure
func (c *Client) FetchFile(fileKey string) (*FigmaFile, error) {
	url := fmt.Sprintf("https://api.figma.com/v1/files/%s", fileKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Figma-Token", c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var file FigmaFile
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &file, nil
}

// FetchComponents extracts all COMPONENT/COMPONENT_SET nodes
func (c *Client) FetchComponents(fileKey string) ([]Component, error) {
	file, err := c.FetchFile(fileKey)
	if err != nil {
		return nil, err
	}

	return TraverseTree(file.Document, fileKey), nil
}

// SearchComponents searches for components matching query
func (c *Client) SearchComponents(fileKey, query string) ([]Component, error) {
	components, err := c.FetchComponents(fileKey)
	if err != nil {
		return nil, err
	}

	return FilterComponents(components, query), nil
}

// GetComponentByNodeID fetches a component by its node ID
func (c *Client) GetComponentByNodeID(fileKey, nodeID string) (*Component, error) {
	components, err := c.FetchComponents(fileKey)
	if err != nil {
		return nil, err
	}

	// Find component by node ID
	for _, comp := range components {
		if comp.NodeID == nodeID {
			return &comp, nil
		}
	}

	return nil, fmt.Errorf("component with node ID %s not found", nodeID)
}

// GetFullMetadata fetches enhanced metadata for a component
func (c *Client) GetFullMetadata(fileKey, componentID string) (*EnhancedMetadata, error) {
	// First get the component
	component, err := c.GetComponentByNodeID(fileKey, componentID)
	if err != nil {
		return nil, err
	}

	// Create enhanced metadata
	metadata := CreateEnhancedMetadata(*component, fileKey)
	return metadata, nil
}

// GetAllComponents fetches all components from a file
func (c *Client) GetAllComponents(fileKey string) ([]Component, error) {
	return c.FetchComponents(fileKey)
}
