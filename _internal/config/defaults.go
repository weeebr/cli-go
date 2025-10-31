package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func (c *Config) SetDefaults() {
	if c.Ringier.DefaultProjectKey == "" {
		c.Ringier.DefaultProjectKey = "PNT"
	}
	if c.History.DefaultDays == 0 {
		c.History.DefaultDays = 7
	}

	if c.AI.Models.OpenAI == "" {
		c.AI.Models.OpenAI = "gpt-4o"
	}
	if c.AI.Models.Anthropic == "" {
		c.AI.Models.Anthropic = "claude-sonnet-4-5-20250929"
	}
	if c.AI.Models.Google == "" {
		c.AI.Models.Google = "gemini-2.5-flash"
	}
	if c.AI.Models.Groq == "" {
		c.AI.Models.Groq = "llama-3.3-70b"
	}
	if c.AI.Timeouts.Default == 0 {
		c.AI.Timeouts.Default = 60
	}

	if c.Network.TimeoutSeconds == 0 {
		c.Network.TimeoutSeconds = 30
	}
	if c.Network.RetryAttempts == 0 {
		c.Network.RetryAttempts = 3
	}

	if c.Prompts.BaseDir == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			c.Prompts.BaseDir = filepath.Join(homeDir, ".prompts")
		} else {
			c.Prompts.BaseDir = ".prompts"
		}
	}

	if c.Cache.BaseDir == "" {
		c.Cache.BaseDir = "_internal/cache"
	}
	if c.Cache.MaxSizeMB == 0 {
		c.Cache.MaxSizeMB = 500
	}
	if c.Cache.TTLDays == 0 {
		c.Cache.TTLDays = 30
	}

	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.ViteURL == "" {
		c.Server.ViteURL = "http://localhost:5173"
	}
	if c.Server.APIPort == 0 {
		c.Server.APIPort = 8888
	}
	if c.Server.FrontendType == "" {
		c.Server.FrontendType = "svelte"
	}
	if c.Server.Environment == "" {
		c.Server.Environment = "development"
	}

	if c.Credentials.File == "" {
		c.Credentials.File = ".cli-go/credentials.enc"
	}
}

func createDefaultConfig() error {
	config := Config{
		Ringier: struct {
			DefaultProjectKey string `json:"defaultProjectKey" yaml:"default_project_key"`
			DefaultUser       string `json:"defaultUser" yaml:"default_user"`
		}{
			DefaultProjectKey: "PROJECT",
			DefaultUser:       "user@example.com",
		},
		Repositories: []RepoConfig{
			// Add your repositories here
			// Example:
			// {
			//     Owner:      "your-org",
			//     Repo:       "your-repo",
			//     Path:       "/path/to/your/repo",
			//     MainBranch: "main",
			//     Main:       true,
			// },
		},
		History: struct {
			DefaultDays  int    `json:"defaultDays" yaml:"default_days"`
			AuthorFilter string `json:"authorFilter" yaml:"author_filter"`
		}{
			DefaultDays:  7,
			AuthorFilter: "Your Name",
		},
		AI: struct {
			Models struct {
				OpenAI    string `json:"openai" yaml:"openai"`
				Anthropic string `json:"anthropic" yaml:"anthropic"`
				Google    string `json:"google" yaml:"google"`
				Groq      string `json:"groq" yaml:"groq"`
			} `json:"models" yaml:"models"`
			Timeouts struct {
				Default int `json:"default" yaml:"default"`
			} `json:"timeouts" yaml:"timeouts"`
		}{
			Models: struct {
				OpenAI    string `json:"openai" yaml:"openai"`
				Anthropic string `json:"anthropic" yaml:"anthropic"`
				Google    string `json:"google" yaml:"google"`
				Groq      string `json:"groq" yaml:"groq"`
			}{
				OpenAI:    "gpt-4o",
				Anthropic: "claude-sonnet-4-5-20250929",
				Google:    "gemini-2.5-flash",
				Groq:      "llama-3.3-70b",
			},
			Timeouts: struct {
				Default int `json:"default" yaml:"default"`
			}{
				Default: 60,
			},
		},
		Network: struct {
			TimeoutSeconds int `json:"timeoutSeconds" yaml:"timeout_seconds"`
			RetryAttempts  int `json:"retryAttempts" yaml:"retry_attempts"`
		}{
			TimeoutSeconds: 30,
			RetryAttempts:  3,
		},
		Prompts: struct {
			BaseDir string `json:"baseDir" yaml:"base_dir"`
		}{
			BaseDir: "~/.prompts",
		},
		Cache: struct {
			BaseDir   string `json:"baseDir" yaml:"base_dir"`
			MaxSizeMB int    `json:"maxSizeMB" yaml:"max_size_mb"`
			TTLDays   int    `json:"ttlDays" yaml:"ttl_days"`
		}{
			BaseDir:   "_internal/cache",
			MaxSizeMB: 500,
			TTLDays:   30,
		},
		Server: struct {
			Port         int    `json:"port" yaml:"port"`
			ViteURL      string `json:"viteURL" yaml:"vite_url"`
			APIPort      int    `json:"apiPort" yaml:"api_port"`
			FrontendType string `json:"frontendType" yaml:"frontend_type"`
			Environment  string `json:"environment" yaml:"environment"`
		}{
			Port:         8080,
			ViteURL:      "http://localhost:5173",
			APIPort:      8888,
			FrontendType: "svelte",
			Environment:  "development",
		},
		Credentials: struct {
			File string `json:"file" yaml:"file"`
		}{
			File: ".cli-go/credentials.enc",
		},
	}

	// Write as YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	return os.WriteFile("config.yml", data, 0644)
}
