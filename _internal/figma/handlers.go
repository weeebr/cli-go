package figma

import (
	"fmt"
	"cli-go/_internal/io"
	"os"
)

// HandleSearch handles search operations
func HandleSearch(args []string, token string, clip bool, file string, compact, json bool) {
	if len(args) == 0 {
		outputError("search", "query required", clip, file, compact, json)
		return
	}

	query := args[0]
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if len(args) > 1 {
		fileKey = args[1]
	}

	// Check if query is a Figma URL
	if IsFigmaURL(query) {
		HandleURLLookup(query, fileKey, clip, file, compact, json)
		return
	}

	// First, search cache for matching components
	cachedComponents, err := SearchCacheByFileKey(query, fileKey)
	if err == nil && len(cachedComponents) > 0 {
		fmt.Fprintf(os.Stderr, "%s\n\n", io.FormatWithEmoji(fmt.Sprintf("Cached: %d matches", len(cachedComponents)), "cache"))
		fmt.Fprint(os.Stderr, FormatResults(cachedComponents, true))
	}

	// Check if we should fetch fresh data
	client := NewClient(token, fileKey)
	if err := client.ValidateToken(); err != nil {
		outputError("search", err.Error(), clip, file, compact, json)
		return
	}

	freshComponents, err := client.SearchComponents(fileKey, query)
	if err != nil {
		outputError("search", fmt.Sprintf("API error: %v", err), clip, file, compact, json)
		return
	}

	// Show fresh results
	fmt.Fprintf(os.Stderr, "%s\n\n", io.FormatWithEmoji(fmt.Sprintf("Fresh: %d matches", len(freshComponents)), "cache"))

	// Format and output results
	formatted := FormatResults(freshComponents, compact)
	io.DirectOutput(formatted, clip, file, false)
}

// HandleURLLookup handles URL lookup operations
func HandleURLLookup(url, fileKey string, clip bool, file string, compact, json bool) {
	// Extract file key and node ID from URL
	extractedFileKey, nodeID, err := ExtractFileAndNodeFromURL(url)
	if err != nil {
		outputError("url", fmt.Sprintf("Invalid URL: %v", err), clip, file, compact, json)
		return
	}

	// Use extracted file key if not provided
	if fileKey == "Bvw817OVY6zhmEty1Syj8Q" {
		fileKey = extractedFileKey
	}

	// Check cache first
	cachedComponent, err := GetCachedComponent(fileKey, nodeID)
	if err == nil && cachedComponent != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", io.FormatWithEmoji("Cached component found", "cache"))
		formatted := FormatComponent(*cachedComponent, compact)
		io.DirectOutput(formatted, clip, file, false)
		return
	}

	// Fetch fresh data
	client := NewClient("", fileKey)
	component, err := client.GetComponentByNodeID(fileKey, nodeID)
	if err != nil {
		outputError("url", fmt.Sprintf("API error: %v", err), clip, file, compact, json)
		return
	}

	// Cache the component
	CacheComponent(*component)

	// Format and output
	formatted := FormatComponent(*component, compact)
	io.DirectOutput(formatted, clip, file, false)
}

// HandleFullMetadata handles full metadata operations
func HandleFullMetadata(args []string, token string, clip bool, file string, compact, json bool) {
	if len(args) == 0 {
		outputError("metadata", "component ID required", clip, file, compact, json)
		return
	}

	componentID := args[0]
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if len(args) > 1 {
		fileKey = args[1]
	}

	// Check cache first
	cachedMetadata, err := GetCachedMetadata(fileKey, componentID)
	if err == nil && cachedMetadata != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", io.FormatWithEmoji("Cached metadata found", "cache"))
		formatted := FormatEnhancedMetadata(cachedMetadata)
		io.DirectOutput(formatted, clip, file, false)
		return
	}

	// Fetch fresh data
	client := NewClient(token, fileKey)
	if err := client.ValidateToken(); err != nil {
		outputError("metadata", err.Error(), clip, file, compact, json)
		return
	}

	metadata, err := client.GetFullMetadata(fileKey, componentID)
	if err != nil {
		outputError("metadata", fmt.Sprintf("API error: %v", err), clip, file, compact, json)
		return
	}

	// Cache the metadata
	CacheMetadata(*metadata)

	// Format and output
	formatted := FormatEnhancedMetadata(metadata)
	io.DirectOutput(formatted, clip, file, false)
}

// HandleInit handles initialization operations
func HandleInit(args []string, token string, clip bool, file string, compact, json bool) {
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if len(args) > 0 {
		fileKey = args[0]
	}

	// Validate token
	client := NewClient(token, fileKey)
	if err := client.ValidateToken(); err != nil {
		outputError("init", err.Error(), clip, file, compact, json)
		return
	}

	// Fetch all components
	components, err := client.GetAllComponents(fileKey)
	if err != nil {
		outputError("init", fmt.Sprintf("API error: %v", err), clip, file, compact, json)
		return
	}

	// Cache all components
	CacheComponents(components)

	// Format and output
	formatted := fmt.Sprintf("Initialized cache with %d components", len(components))
	io.DirectOutput(formatted, clip, file, false)
}

// HandleList handles list operations
func HandleList(args []string, token string, clip bool, file string, compact, json bool) {
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if len(args) > 0 {
		fileKey = args[0]
	}

	// Get cached components
	components, err := GetCachedComponents(fileKey)
	if err != nil {
		outputError("list", fmt.Sprintf("Cache error: %v", err), clip, file, compact, json)
		return
	}

	if len(components) == 0 {
		outputError("list", "No cached components found. Run 'figma init' first.", clip, file, compact, json)
		return
	}

	// Format and output
	formatted := FormatResults(components, compact)
	io.DirectOutput(formatted, clip, file, false)
}

// HandleCache handles cache operations
func HandleCache(args []string, token string, clip bool, file string, compact, json bool) {
	if len(args) == 0 {
		outputError("cache", "operation required (clear|stats|info)", clip, file, compact, json)
		return
	}

	operation := args[0]
	fileKey := "Bvw817OVY6zhmEty1Syj8Q" // Default file key
	if len(args) > 1 {
		fileKey = args[1]
	}

	switch operation {
	case "clear":
		err := ClearCache(fileKey)
		if err != nil {
			outputError("cache", fmt.Sprintf("Clear error: %v", err), clip, file, compact, json)
			return
		}
		formatted := "Cache cleared successfully"
		io.DirectOutput(formatted, clip, file, false)

	case "stats":
		entries, lastModified, err := GetCacheStats(fileKey)
		if err != nil {
			outputError("cache", fmt.Sprintf("Stats error: %v", err), clip, file, compact, json)
			return
		}
		formatted := FormatCacheStats(entries, lastModified)
		io.DirectOutput(formatted, clip, file, false)

	case "info":
		info, err := GetCacheInfo(fileKey)
		if err != nil {
			outputError("cache", fmt.Sprintf("Info error: %v", err), clip, file, compact, json)
			return
		}
		formatted := FormatCacheInfo(info)
		io.DirectOutput(formatted, clip, file, false)

	default:
		outputError("cache", "Invalid operation. Use: clear, stats, or info", clip, file, compact, json)
	}
}

// outputError outputs an error message
func outputError(action, message string, clip bool, file string, compact, json bool) {
	errorMsg := fmt.Sprintf("Error in %s: %s", action, message)
	io.DirectOutput(errorMsg, clip, file, false)
}
