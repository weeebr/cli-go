package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Set stores a key-value pair with tags
func (s *Store) Set(key string, data map[string]interface{}, tags []string) error {
	entry := Entry{
		Key:       key,
		Data:      data,
		Tags:      tags,
		Timestamp: time.Now(),
	}

	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %v", err)
	}

	filename := fmt.Sprintf("%s.json", key)
	filepath := filepath.Join(s.dir, filename)

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %v", err)
	}

	return nil
}

// Get retrieves an entry by key
func (s *Store) Get(key string) (*Entry, error) {
	filename := fmt.Sprintf("%s.json", key)
	filepath := filepath.Join(s.dir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("entry not found")
		}
		return nil, fmt.Errorf("failed to read cache file: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry: %v", err)
	}

	return &entry, nil
}

// Search finds entries with similar tags using Jaccard similarity
func (s *Store) Search(tags []string, threshold float64) ([]*Entry, error) {
	if len(tags) == 0 {
		return []*Entry{}, nil
	}

	entries, err := s.GetAllEntries()
	if err != nil {
		return nil, err
	}

	var results []*Entry
	for _, entry := range entries {
		similarity := jaccardSimilarity(tags, entry.Tags)
		if similarity >= threshold {
			results = append(results, entry)
		}
	}

	// Sort by similarity (descending) and timestamp (descending)
	sort.Slice(results, func(i, j int) bool {
		simI := jaccardSimilarity(tags, results[i].Tags)
		simJ := jaccardSimilarity(tags, results[j].Tags)
		if simI == simJ {
			return results[i].Timestamp.After(results[j].Timestamp)
		}
		return simI > simJ
	})

	return results, nil
}

// Update updates an existing entry
func (s *Store) Update(key string, data map[string]interface{}, tags []string) error {
	// Check if entry exists
	_, err := s.Get(key)
	if err != nil {
		return fmt.Errorf("entry not found: %v", err)
	}

	// Update the entry
	return s.Set(key, data, tags)
}

// Clear removes all entries from the cache
func (s *Store) Clear() error {
	pattern := filepath.Join(s.dir, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list cache files: %v", err)
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove cache file: %v", err)
		}
	}

	return nil
}

// Delete removes a specific entry by key
func (s *Store) Delete(key string) error {
	filename := fmt.Sprintf("%s.json", key)
	filepath := filepath.Join(s.dir, filename)

	if err := os.Remove(filepath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("entry not found")
		}
		return fmt.Errorf("failed to remove cache file: %v", err)
	}

	return nil
}

// GetAllEntries retrieves all entries from the cache directory
func (s *Store) GetAllEntries() ([]*Entry, error) {
	pattern := filepath.Join(s.dir, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list cache files: %v", err)
	}

	var entries []*Entry
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip corrupted files
		}

		var entry Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue // Skip invalid JSON
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

// jaccardSimilarity calculates Jaccard similarity between two tag sets
func jaccardSimilarity(tags1, tags2 []string) float64 {
	if len(tags1) == 0 && len(tags2) == 0 {
		return 1.0
	}
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0.0
	}

	set1 := make(map[string]bool)
	for _, tag := range tags1 {
		set1[tag] = true
	}

	set2 := make(map[string]bool)
	for _, tag := range tags2 {
		set2[tag] = true
	}

	intersection := 0
	for tag := range set1 {
		if set2[tag] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// ExtractTags extracts meaningful tags from text
func ExtractTags(text string) []string {
	words := strings.Fields(strings.ToLower(text))

	// Filter words longer than 2 characters
	var tags []string
	seen := make(map[string]bool)

	for _, word := range words {
		// Remove common punctuation and special characters
		word = strings.Trim(word, ".,!?;:\"'()[]{}@#$%^&*_+{}|:<>?[]\\;'\",./")
		// Only keep alphanumeric words
		if len(word) > 2 && isAlphanumeric(word) && !seen[word] {
			tags = append(tags, word)
			seen[word] = true
		}
	}

	return tags
}

// isAlphanumeric checks if a string contains only alphanumeric characters
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// GetStats returns cache statistics
func (s *Store) GetStats() (*StoreStats, error) {
	entries, err := s.GetAllEntries()
	if err != nil {
		return nil, err
	}

	var totalSize int64
	var lastModified time.Time

	for _, entry := range entries {
		// Calculate approximate size
		jsonData, _ := json.Marshal(entry)
		totalSize += int64(len(jsonData))

		if entry.Timestamp.After(lastModified) {
			lastModified = entry.Timestamp
		}
	}

	return &StoreStats{
		Entries:      len(entries),
		TotalSize:    totalSize,
		Namespace:    s.namespace,
		CacheDir:     s.dir,
		LastModified: lastModified,
	}, nil
}

// SetWithOperation stores a key-value pair and returns operation details
func (s *Store) SetWithOperation(key string, data map[string]interface{}, tags []string) (*CacheOperation, error) {
	// Check if entry exists before storing
	_, err := s.Get(key)
	isUpdate := (err == nil)
	
	// Store the entry
	err = s.Set(key, data, tags)
	if err != nil {
		return nil, err
	}
	
	if isUpdate {
		return &CacheOperation{
			Updated: 1,
			Total:   1,
		}, nil
	}
	
	return &CacheOperation{
		Added: 1,
		Total: 1,
	}, nil
}

// ClearWithOperation clears all entries and returns operation details
func (s *Store) ClearWithOperation() (*CacheOperation, error) {
	entries, err := s.List()
	if err != nil {
		return nil, err
	}

	count := len(entries)
	for _, entry := range entries {
		if err := s.Delete(entry.Key); err != nil {
			return nil, err
		}
	}

	return &CacheOperation{
		Removed: count,
		Total:   count,
	}, nil
}

// GenerateKey creates a cache key from namespace and identifier
func GenerateKey(namespace, identifier string) string {
	return fmt.Sprintf("%s:%s", namespace, identifier)
}

// List returns all entries in the cache
func (s *Store) List() ([]Entry, error) {
	entries := []Entry{}

	files, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filepath := filepath.Join(s.dir, file.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			continue
		}

		var entry Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
