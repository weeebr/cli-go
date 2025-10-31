package figma

import (
	"fmt"
	"cli-go/_internal/cache"
	"os"
	"strings"
)

// CacheStore is the Figma components cache store
var CacheStore *cache.Store

func init() {
	var err error
	CacheStore, err = cache.New("figma")
	if err != nil {
		// Log error but don't panic - will be handled in functions
		fmt.Printf("Warning: failed to initialize cache store: %v\n", err)
	}
}

// LoadCache loads all cached components
func LoadCache() (map[string]cache.Entry, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	// Use Search to get all figma entries
	entries, err := CacheStore.Search([]string{"figma"}, 0.0)
	if err != nil {
		return nil, err
	}

	result := make(map[string]cache.Entry)
	for _, entry := range entries {
		result[entry.Key] = *entry
	}
	return result, nil
}

// SaveCache saves components to cache
func SaveCache(components []Component, fileKey string) error {
	if CacheStore == nil {
		return fmt.Errorf("cache store not initialized")
	}

	for _, comp := range components {
		key := "figma-" + fileKey + "-" + comp.NodeID
		data := map[string]interface{}{
			"node_id": comp.NodeID,
			"name":    comp.Name,
			"type":    comp.Type,
			"url":     comp.URL,
		}
		tags := []string{"figma", "component", fileKey}

		if err := CacheStore.Set(key, data, tags); err != nil {
			return err
		}
	}
	return nil
}

// SearchCacheByFileKey searches cache for components matching query in specific file
func SearchCacheByFileKey(query, fileKey string) ([]Component, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	// Search for entries with both figma and fileKey tags
	entries, err := CacheStore.Search([]string{"figma", fileKey}, 0.0)
	if err != nil {
		return nil, err
	}

	var results []Component
	query = strings.ToLower(query)

	for _, entry := range entries {
		// Convert Data to Component using safe extraction
		comp, err := ComponentFromCacheData(entry.Data)
		if err != nil {
			continue // Skip invalid entries
		}
		if strings.Contains(strings.ToLower(comp.Name), query) {
			results = append(results, comp)
		}
	}

	return results, nil
}

// SearchCache searches cache for components matching query
func SearchCache(query string) ([]Component, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	entries, err := CacheStore.Search([]string{"figma"}, 0.0)
	if err != nil {
		return nil, err
	}

	var results []Component
	query = strings.ToLower(query)

	for _, entry := range entries {
		// Convert Data to Component using safe extraction
		comp, err := ComponentFromCacheData(entry.Data)
		if err != nil {
			continue // Skip invalid entries
		}
		if strings.Contains(strings.ToLower(comp.Name), query) {
			results = append(results, comp)
		}
	}

	return results, nil
}

// ClearCache removes all cached components for a specific file
func ClearCache(fileKey string) error {
	if CacheStore == nil {
		return fmt.Errorf("cache store not initialized")
	}

	// Get all entries for this file
	entries, err := CacheStore.Search([]string{"figma", fileKey}, 0.0)
	if err != nil {
		return err
	}

	// Delete each entry
	for _, entry := range entries {
		if err := CacheStore.Delete(entry.Key); err != nil {
			return fmt.Errorf("failed to delete cache entry %s: %v", entry.Key, err)
		}
	}

	return nil
}

// GetCachedComponent retrieves a single component from cache
func GetCachedComponent(fileKey, nodeID string) (*Component, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	key := "figma-" + fileKey + "-" + nodeID
	entry, err := CacheStore.Get(key)
	if err != nil {
		return nil, err
	}

	comp, err := ComponentFromCacheData(entry.Data)
	if err != nil {
		return nil, err
	}

	return &comp, nil
}

// GetCachedMetadata retrieves metadata from cache
func GetCachedMetadata(fileKey, componentID string) (*EnhancedMetadata, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	key := "figma-metadata-" + fileKey + "-" + componentID
	entry, err := CacheStore.Get(key)
	if err != nil {
		return nil, err
	}

	// Convert cache data to EnhancedMetadata
	metadata := &EnhancedMetadata{}
	if nodeID, ok := entry.Data["node_id"].(string); ok {
		metadata.Component.NodeID = nodeID
	}
	if name, ok := entry.Data["name"].(string); ok {
		metadata.Component.Name = name
	}
	if compType, ok := entry.Data["type"].(string); ok {
		metadata.Component.Type = compType
	}
	if url, ok := entry.Data["url"].(string); ok {
		metadata.Component.URL = url
	}

	return metadata, nil
}

// CacheComponent caches a single component
func CacheComponent(component Component) error {
	if CacheStore == nil {
		return fmt.Errorf("cache store not initialized")
	}

	// Extract file key from component URL or use default
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if component.URL != "" {
		// Try to extract file key from URL
		if extracted, err := ExtractFileIDFromURL(component.URL); err == nil {
			fileKey = extracted
		}
	}

	key := "figma-" + fileKey + "-" + component.NodeID
	data := map[string]interface{}{
		"node_id": component.NodeID,
		"name":    component.Name,
		"type":    component.Type,
		"url":     component.URL,
	}
	tags := []string{"figma", "component", fileKey}

	return CacheStore.Set(key, data, tags)
}

// CacheMetadata caches enhanced metadata
func CacheMetadata(metadata EnhancedMetadata) error {
	if CacheStore == nil {
		return fmt.Errorf("cache store not initialized")
	}

	// Extract file key from component URL or use default
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if metadata.Component.URL != "" {
		// Try to extract file key from URL
		if extracted, err := ExtractFileIDFromURL(metadata.Component.URL); err == nil {
			fileKey = extracted
		}
	}

	key := "figma-metadata-" + fileKey + "-" + metadata.Component.NodeID
	data := map[string]interface{}{
		"node_id":          metadata.Component.NodeID,
		"name":             metadata.Component.Name,
		"type":             metadata.Component.Type,
		"url":              metadata.Component.URL,
		"description":      metadata.Description,
		"has_interactions": metadata.HasInteractions,
		"property_count":   metadata.PropertyCount,
	}
	tags := []string{"figma", "metadata", fileKey}

	return CacheStore.Set(key, data, tags)
}

// GetCachedComponents retrieves all cached components for a file
func GetCachedComponents(fileKey string) ([]Component, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	entries, err := CacheStore.Search([]string{"figma", fileKey}, 0.0)
	if err != nil {
		return nil, err
	}

	var components []Component
	for _, entry := range entries {
		// Skip metadata entries
		if strings.Contains(entry.Key, "metadata") {
			continue
		}

		comp, err := ComponentFromCacheData(entry.Data)
		if err != nil {
			continue // Skip invalid entries
		}
		components = append(components, comp)
	}

	return components, nil
}

// CacheComponents caches multiple components
func CacheComponents(components []Component) error {
	if CacheStore == nil {
		return fmt.Errorf("cache store not initialized")
	}

	// Extract file key from first component URL or use default
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if len(components) > 0 && components[0].URL != "" {
		// Try to extract file key from URL
		if extracted, err := ExtractFileIDFromURL(components[0].URL); err == nil {
			fileKey = extracted
		}
	}

	return SaveCache(components, fileKey)
}

// GetCacheInfo returns detailed cache information
func GetCacheInfo(fileKey string) (map[string]interface{}, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	entries, err := CacheStore.Search([]string{"figma", fileKey}, 0.0)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"file_key":   fileKey,
		"entries":    len(entries),
		"components": 0,
		"metadata":   0,
	}

	for _, entry := range entries {
		if strings.Contains(entry.Key, "metadata") {
			info["metadata"] = info["metadata"].(int) + 1
		} else {
			info["components"] = info["components"].(int) + 1
		}
	}

	return info, nil
}

// GetCacheStats returns cache statistics for a specific file
func GetCacheStats(fileKey string) (int, int64, error) {
	if CacheStore == nil {
		return 0, 0, fmt.Errorf("cache store not initialized")
	}

	// Get entries for this specific file
	entries, err := CacheStore.Search([]string{"figma", fileKey}, 0.0)
	if err != nil {
		return 0, 0, err
	}

	// Get overall stats for timestamp
	stats, err := CacheStore.GetStats()
	if err != nil {
		return 0, 0, err
	}

	return len(entries), stats.LastModified.Unix(), nil
}

// GetCacheStatsDetailed returns detailed cache statistics
func GetCacheStatsDetailed() (*cache.StoreStats, error) {
	if CacheStore == nil {
		return nil, fmt.Errorf("cache store not initialized")
	}

	return CacheStore.GetStats()
}

// UpdateCache merges new components into existing cache
func UpdateCache(components []Component, fileKey string) error {
	// Add new components
	for _, comp := range components {
		key := "figma-" + fileKey + "-" + comp.NodeID
		data := map[string]interface{}{
			"node_id": comp.NodeID,
			"name":    comp.Name,
			"type":    comp.Type,
			"url":     comp.URL,
		}
		tags := []string{"figma", "component", fileKey}

		if err := CacheStore.Set(key, data, tags); err != nil {
			return err
		}
	}
	return nil
}

// UpdateCacheWithDeletions updates cache with fresh components and removes deleted ones
func UpdateCacheWithDeletions(freshComponents []Component, fileKey string) error {
	if CacheStore == nil {
		return fmt.Errorf("cache store not initialized")
	}

	// Get all existing cached components for this file
	existingEntries, err := CacheStore.Search([]string{"figma", fileKey}, 0.0)
	if err != nil {
		return fmt.Errorf("failed to load existing cache: %v", err)
	}

	// Create a map of fresh component node IDs for quick lookup
	freshNodeIDs := make(map[string]bool)
	for _, comp := range freshComponents {
		freshNodeIDs[comp.NodeID] = true
	}

	// Remove components that are no longer in the fresh results
	for _, entry := range existingEntries {
		nodeID, ok := entry.Data["node_id"].(string)
		if !ok {
			continue // Skip invalid entries
		}

		if !freshNodeIDs[nodeID] {
			// Component no longer exists, remove from cache
			if err := CacheStore.Delete(entry.Key); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to delete cache entry %s: %v\n", entry.Key, err)
			}
		}
	}

	// Add/update fresh components
	for _, comp := range freshComponents {
		key := "figma-" + fileKey + "-" + comp.NodeID
		data := map[string]interface{}{
			"node_id": comp.NodeID,
			"name":    comp.Name,
			"type":    comp.Type,
			"url":     comp.URL,
		}
		tags := []string{"figma", "component", fileKey}

		if err := CacheStore.Set(key, data, tags); err != nil {
			return fmt.Errorf("failed to update cache for component %s: %v", comp.NodeID, err)
		}
	}

	return nil
}
