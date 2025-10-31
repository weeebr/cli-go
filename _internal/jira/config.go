package jira

import (
	"fmt"
	"cli-go/_internal/config"
)

// Config represents Jira configuration from config.yml
type Config struct {
	Email          string
	DefaultProject string
	BaseURL        string
}

// LoadConfig loads Jira configuration from the unified config
func LoadConfig() (*Config, error) {
	// Load the unified configuration
	unifiedConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Extract Jira configuration
	jiraConfig := &Config{
		BaseURL:        unifiedConfig.Jira.BaseURL,
		Email:          unifiedConfig.Jira.Email,
		DefaultProject: unifiedConfig.Jira.DefaultProject,
	}

	// Validate required fields
	if jiraConfig.BaseURL == "" {
		return nil, fmt.Errorf("jira.base_url not configured")
	}
	if jiraConfig.Email == "" {
		return nil, fmt.Errorf("jira.email not configured")
	}
	if jiraConfig.DefaultProject == "" {
		return nil, fmt.Errorf("jira.default_project not configured")
	}

	return jiraConfig, nil
}

// LoadJiraConfig loads Jira configuration and API token
func LoadJiraConfig() (*Config, string, error) {
	// Load configuration
	jiraConfig, err := LoadConfig()
	if err != nil {
		return nil, "", err
	}

	// Load API token from credentials
	apiToken, err := config.GetKey("jira")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get Jira API token: %v", err)
	}

	return jiraConfig, apiToken, nil
}
