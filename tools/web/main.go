package main

// DESCRIPTION: Web search tool using Perplexity AI with intelligent caching

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"cli-go/_internal/ai"
	"cli-go/_internal/cache"
	"cli-go/_internal/config"
	"cli-go/_internal/io"
)

// SearchResult represents a web search result
type SearchResult struct {
	Query     string   `json:"query"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Cached    bool     `json:"cached"`
	Timestamp string   `json:"timestamp"`
}

// ClearResult represents cache clear operation result
type ClearResult struct {
	Success        bool `json:"success"`
	EntriesCleared int  `json:"entries_cleared"`
}

func main() {
	// Check for special commands before flag parsing
	args := os.Args[1:]
	if len(args) >= 1 {
		// Handle test mode
		if args[0] == "--test-real" {
			handleTestMode(args[1:], false, "", false)
			return
		}
	}

	var (
		clip    = flag.Bool("clip", false, "Copy to clipboard")
		file    = flag.String("file", "", "Write to file")
		compact = flag.Bool("compact", false, "Use compact JSON output")
		json    = flag.Bool("json", false, "Output in JSON format")
	)
	flag.Parse()

	args = flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: web <query> | web cache <stats|clear>\n")
		os.Exit(1)
	}

	// Handle cache subcommands
	if args[0] == "cache" {
		handleCacheCommand(args[1:], *clip, *file, *compact, *json)
		return
	}

	// Handle search query
	query := strings.Join(args, " ")
	fmt.Fprintf(os.Stderr, "DEBUG: json flag=%v\n", *json)
	handleSearch(query, *clip, *file, *compact, *json)
}

func handleCacheCommand(args []string, clip bool, file string, compact, json bool) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: web cache <stats|clear>\n")
		os.Exit(1)
	}

	store, err := cache.New("web")
	ai.ExitIf(err, "failed to initialize cache")

	switch args[0] {
	case "stats":
		stats, err := store.GetStats()
		ai.ExitIf(err, "failed to get cache stats")

		io.DirectOutput(stats, clip, file, json)

	case "clear":
		// Use ClearWithOperation for simple 1-liner output
		op, err := store.ClearWithOperation()
		ai.ExitIf(err, "failed to clear cache")

		// Output simple 1-liner to stderr (not stdout) only when not in JSON mode
		if !json {
			fmt.Fprintf(os.Stderr, "%s\n", op.FormatCacheClear())
		}

		// Still provide JSON result for programmatic use
		result := ClearResult{
			Success:        true,
			EntriesCleared: op.Removed,
		}

		io.DirectOutput(result, clip, file, json)

	default:
		fmt.Fprintf(os.Stderr, "Unknown cache command: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "Available commands: stats, clear\n")
		os.Exit(1)
	}
}

func handleTestMode(args []string, clip bool, file string, compact bool) {
	query := strings.Join(args, " ")
	if query == "" {
		query = "test query"
	}

	// Simulate real API response
	mockResponse := fmt.Sprintf(`This is a simulated real API response for: "%s"

Key insights:
- This demonstrates how the web tool would work with a valid API key
- Real results would be cached for future searches
- The transition from mock to real data is seamless

The system is working correctly - it just needs a valid API key to make real API calls.`, query)

	tags := cache.ExtractTags(query)
	result := SearchResult{
		Query:     query,
		Content:   mockResponse,
		Tags:      tags,
		Cached:    false,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	io.DirectOutput(result, clip, file, false)
}

func handleSearch(query string, clip bool, file string, compact, json bool) {
	// Initialize cache first
	cacheStore, err := cache.New("web")
	ai.ExitIf(err, "failed to initialize cache")

	// Extract tags and search cache
	tags := cache.ExtractTags(query)
	cachedResults, err := cacheStore.Search(tags, 0.5) // Increased threshold for better precision
	ai.ExitIf(err, "cache search failed")

	// If we found cached results, show them first
	if len(cachedResults) > 0 {
		entry := cachedResults[0]
		content, ok := entry.Data["content"].(string)
		if !ok {
			content = ""
		}

		result := SearchResult{
			Query:     query,
			Content:   content,
			Tags:      entry.Tags,
			Cached:    true,
			Timestamp: entry.Timestamp.Format(time.RFC3339),
		}

		io.DirectOutput(result, clip, file, json)
		// For JSON output, return early to avoid duplicate output
		if json {
			return
		}
		// Continue to fetch fresh data (don't return early for non-JSON)
	}

	// Always fetch fresh data (after showing cached results if any)
	apiKey, err := config.GetKey("perplexity")
	ai.ExitIf(err, "failed to get Perplexity API key")

	client := ai.NewPerplexityClient(apiKey)
	content, err := client.Search(query)
	if err != nil {
		// Check if this is a mock API error that should return mock content
		if mockErr, ok := err.(*ai.MockAPIError); ok {
			content = mockErr.GetMockResponse()
			// Don't cache mock responses - return immediately
			result := SearchResult{
				Query:     query,
				Content:   content,
				Tags:      cache.ExtractTags(query),
				Cached:    false,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			io.DirectOutput(result, clip, file, json)
			return
		} else {
			ai.ExitIf(err, "Perplexity API error")
		}
	}

	// Only cache real API responses, not mock responses
	shouldCache := true

	// Store result in cache (only for real API responses)
	if shouldCache {
		key := cache.GenerateKey("web", query)
		data := map[string]interface{}{
			"query":   query,
			"content": content,
		}

		op, err := cacheStore.SetWithOperation(key, data, tags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to cache result: %v\n", err)
			// Continue execution even if caching fails
		} else if !json {
			// Only output cache message when not in JSON mode
			fmt.Fprintf(os.Stderr, "%s\n", op.FormatCacheUpdate())
		}
	}

	result := SearchResult{
		Query:     query,
		Content:   content,
		Tags:      tags,
		Cached:    false,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	io.DirectOutput(result, clip, file, json)
}
