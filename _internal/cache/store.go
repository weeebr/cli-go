package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cli-go/_internal/config"
)

// Store represents a generic cache store with namespace isolation
type Store struct {
	dir       string
	namespace string
}

// CacheOperation represents the result of a cache operation
type CacheOperation struct {
	Added   int
	Removed int
	Updated int
	Total   int
}

// FormatCacheClear returns a formatted string for cache clear operations
func (op *CacheOperation) FormatCacheClear() string {
	if op.Removed > 0 {
		return fmt.Sprintf("Cleared %d cache entries", op.Removed)
	}
	return "No cache entries to clear"
}

// FormatCacheUpdate returns a formatted string for cache update operations
func (op *CacheOperation) FormatCacheUpdate() string {
	if op.Added > 0 {
		return fmt.Sprintf("Added %d cache entries", op.Added)
	}
	return "No cache entries added"
}

// Entry represents a cached entry with metadata
type Entry struct {
	Key       string                 `json:"key"`
	Data      map[string]interface{} `json:"data"`
	Tags      []string               `json:"tags"`
	Timestamp time.Time              `json:"timestamp"`
}

// StoreStats represents cache statistics
type StoreStats struct {
	Entries      int       `json:"entries"`
	TotalSize    int64     `json:"total_size_bytes"`
	Namespace    string    `json:"namespace"`
	CacheDir     string    `json:"cache_dir"`
	LastModified time.Time `json:"last_modified,omitempty"`
}

// findProjectRoot finds the project root by walking up to find go.mod
func findProjectRoot() (string, error) {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

// New creates a new cache store for the given namespace
func New(namespace string) (*Store, error) {
	// Load config to get cache directory
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fallback to project root cache
		projectRoot, err := findProjectRoot()
		if err != nil {
			return nil, fmt.Errorf("failed to find project root: %v", err)
		}
		dir := filepath.Join(projectRoot, ".cache", namespace)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %v", err)
		}
		return &Store{dir: dir, namespace: namespace}, nil
	}

	// Use config cache directory
	baseDir := cfg.Cache.BaseDir
	if strings.HasPrefix(baseDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			baseDir = strings.Replace(baseDir, "~", homeDir, 1)
		}
	}

	dir := filepath.Join(baseDir, namespace)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %v", err)
	}

	return &Store{
		dir:       dir,
		namespace: namespace,
	}, nil
}
