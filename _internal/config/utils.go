package config

import (
	"fmt"
	"os"
)

func (c *Config) GetMainRepos() []RepoConfig {
	var main []RepoConfig
	for _, repo := range c.Repositories {
		if repo.Main {
			main = append(main, repo)
		}
	}
	return main
}

func (c *Config) GetRepoByPath(path string) *RepoConfig {
	for _, repo := range c.Repositories {
		if repo.Path == path {
			return &repo
		}
	}
	return nil
}

func (c *Config) GetRepoPaths() []string {
	var paths []string
	for _, repo := range c.Repositories {
		paths = append(paths, repo.Path)
	}
	return paths
}

func (c *Config) GetMainRepoPaths() []string {
	var paths []string
	for _, repo := range c.Repositories {
		if repo.Main {
			paths = append(paths, repo.Path)
		}
	}
	return paths
}

func (c *Config) NormalizeTicketID(input string) string {
	if len(input) > 4 && input[3] == '-' {
		return input
	}
	if len(input) > 0 && input[0] >= '0' && input[0] <= '9' {
		return fmt.Sprintf("%s-%s", c.Ringier.DefaultProjectKey, input)
	}
	return input
}

func (c *Config) ValidateConfig() error {
	if len(c.Repositories) == 0 {
		return fmt.Errorf("no repositories configured")
	}
	if c.Ringier.DefaultUser == "" {
		return fmt.Errorf("default user not configured")
	}
	for _, repo := range c.Repositories {
		if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
			return fmt.Errorf("repository path does not exist: %s", repo.Path)
		}
	}
	return nil
}
