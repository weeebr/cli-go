package git

import (
	"fmt"
	"cli-go/_internal/config"
)

// GetReposToProcess determines which repositories to process based on flags
func GetReposToProcess(singlePath string, mainFlag bool, allFlag bool, defaultMode string) ([]string, error) {
	// If explicit single path provided, validate and return it
	if singlePath != "" {
		if !IsGitRepoAtPath(singlePath) {
			return nil, fmt.Errorf("path is not a git repository: %s", singlePath)
		}
		return []string{singlePath}, nil
	}

	// If all flag specified, return all repos
	if allFlag {
		configData, err := config.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %v", err)
		}
		return configData.GetRepoPaths(), nil
	}

	// If main flag specified, return main repos
	if mainFlag {
		configData, err := config.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %v", err)
		}
		return configData.GetMainRepoPaths(), nil
	}

	// No flags specified - use default mode
	switch defaultMode {
	case "pwd":
		// Current repository
		repoPath, err := GetGitRoot()
		if err != nil {
			return nil, fmt.Errorf("not in a git repository: %v", err)
		}
		return []string{repoPath}, nil
	case "main":
		// Main repositories
		configData, err := config.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %v", err)
		}
		return configData.GetMainRepoPaths(), nil
	case "all":
		// All repositories
		configData, err := config.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %v", err)
		}
		return configData.GetRepoPaths(), nil
	default:
		return nil, fmt.Errorf("invalid default mode: %s", defaultMode)
	}
}

// RunAcrossRepos executes an operation across multiple repositories
func RunAcrossRepos(repoPaths []string, operation func(string) (interface{}, error)) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for _, repoPath := range repoPaths {
		result, err := operation(repoPath)
		if err != nil {
			// Log warning but continue with other repos
			results[repoPath] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}
		results[repoPath] = result
	}

	return results, nil
}
