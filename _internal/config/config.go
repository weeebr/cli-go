package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func LoadConfig() (*Config, error) {
	return loadConfig()
}

func CreateDefaultConfig() error {
	return createDefaultConfig()
}

func ValidateConfig(cfg *Config) error {
	return cfg.ValidateConfig()
}

func GetConfigPath() (string, error) {
	// Get the directory of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	
	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)
	
	// Look for config.yml in the cli-go directory (parent of bin/)
	goCliDir := filepath.Join(execDir, "..")
	return filepath.Join(goCliDir, "config.yml"), nil
}

func ConfigExists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func GetConfigInfo() (map[string]interface{}, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	configPath, _ := GetConfigPath()

	return map[string]interface{}{
		"config_file":     configPath,
		"repositories":    len(cfg.Repositories),
		"main_repos":      len(cfg.GetMainRepos()),
		"jira_configured": cfg.Jira.BaseURL != "",
		"ai_models": map[string]string{
			"openai":    cfg.AI.Models.OpenAI,
			"anthropic": cfg.AI.Models.Anthropic,
			"google":    cfg.AI.Models.Google,
			"groq":      cfg.AI.Models.Groq,
		},
		"cache_dir":   cfg.Cache.BaseDir,
		"prompts_dir": cfg.Prompts.BaseDir,
		"server_port": cfg.Server.Port,
	}, nil
}
