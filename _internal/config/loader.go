package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func loadConfig() (*Config, error) {
	// Get the directory of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	
	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)
	
	// Look for config.yml in the cli-go directory (parent of bin/)
	goCliDir := filepath.Join(execDir, "..")
	configPath := filepath.Join(goCliDir, "config.yml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.yml not found in cli-go directory: %s", goCliDir)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yml: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config.yml as YAML or JSON: %v", err)
		}
	}

	config.SetDefaults()

	return &config, nil
}

