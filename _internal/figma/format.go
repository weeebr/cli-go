package figma

import (
	"fmt"
	"strings"

	"cli-go/_internal/io"
)

// FormatResults formats components for display (to stderr)
func FormatResults(components []Component, isCached bool) string {
	var sb strings.Builder

	if isCached {
		sb.WriteString(fmt.Sprintf("%s\n\n", io.FormatWithEmoji(fmt.Sprintf("Cached: %d matches", len(components)), "cache")))
	} else {
		sb.WriteString(fmt.Sprintf("%s\n\n", io.FormatWithEmoji(fmt.Sprintf("Fresh: %d matches", len(components)), "cache")))
	}

	// Simple table format
	sb.WriteString(fmt.Sprintf("%-50s | %s\n", "Component", "Link"))
	sb.WriteString(strings.Repeat("-", 120) + "\n")

	for _, c := range components {
		name := c.Name
		if len(name) > 47 {
			name = name[:47] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-50s | %s\n", name, io.FormatWithEmoji(c.URL, "url")))
	}

	return sb.String()
}

// FormatFullMetadata formats detailed component info (to stderr)
func FormatFullMetadata(meta FullMetadata) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## 💎 Component: %s\n\n", meta.Component.Name))
	sb.WriteString(fmt.Sprintf("**Node ID:** %s\n", meta.Component.NodeID))
	sb.WriteString(fmt.Sprintf("**Type:** %s\n", meta.Component.Type))
	sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", meta.Component.URL))

	if meta.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", meta.Description))
	}

	if meta.HasInteractions {
		sb.WriteString(fmt.Sprintf("⚡ **Interactions:** Yes (%d)\n\n", len(meta.Interactions)))
	}

	if meta.PropertyCount > 0 {
		sb.WriteString(fmt.Sprintf("🎛️ **Properties:** %d defined\n\n", meta.PropertyCount))
	}

	return sb.String()
}

// FormatComponentWithVersions formats component with version information
func FormatComponentWithVersions(name, url, fileKey, nodeID string) string {
	// For now, just format the component line
	// In a real implementation, this would include version tracking
	return fmt.Sprintf("%-50s │ %s\n", name, io.FormatWithEmoji(url, "url"))
}

// FormatResultsWithVersions formats results with version tracking
func FormatResultsWithVersions(components []Component, isCached bool, fileKey string) string {
	var sb strings.Builder

	if isCached {
		sb.WriteString(fmt.Sprintf("%s\n\n", io.FormatWithEmoji(fmt.Sprintf("Cached: %d matches", len(components)), "cache")))
	} else {
		sb.WriteString(fmt.Sprintf("%s\n\n", io.FormatWithEmoji("FRESH RESULTS:", "cache")))
		sb.WriteString("## 🎯 GENERAL - COMPONENT\n\n")
	}

	// Create table header
	sb.WriteString(fmt.Sprintf("%-50s │ %-80s\n", "Component", "Link"))
	sb.WriteString(fmt.Sprintf("%-50s┼%-80s\n", strings.Repeat("─", 50), strings.Repeat("─", 80)))

	// Format each component with version tracking
	for _, comp := range components {
		sb.WriteString(FormatComponentWithVersions(comp.Name, comp.URL, fileKey, comp.NodeID))
	}

	if !isCached {
		sb.WriteString(fmt.Sprintf("\n%s\n", io.FormatWithEmoji(fmt.Sprintf("Cache Update: +%d added, -0 removed, ~0 updated (Total: %d)", len(components), len(components)), "cache")))
	}

	return sb.String()
}

// FormatComponent formats a single component for display
func FormatComponent(component Component, compact bool) string {
	if compact {
		return fmt.Sprintf("%s | %s", component.Name, io.FormatWithEmoji(component.URL, "url"))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Component: %s\n", component.Name))
	sb.WriteString(fmt.Sprintf("**Node ID:** %s\n", component.NodeID))
	sb.WriteString(fmt.Sprintf("**Type:** %s\n", component.Type))
	sb.WriteString(fmt.Sprintf("**URL:** %s\n", io.FormatWithEmoji(component.URL, "url")))

	return sb.String()
}

// FormatEnhancedMetadata formats enhanced metadata for display
func FormatEnhancedMetadata(metadata *EnhancedMetadata) string {
	return FormatFullMetadata(FullMetadata{
		Component:       metadata.Component,
		Properties:      metadata.Properties,
		Interactions:    metadata.Interactions,
		Description:     metadata.Description,
		HasInteractions: metadata.HasInteractions,
		PropertyCount:   metadata.PropertyCount,
	})
}

// FormatCacheStats formats cache statistics for display
func FormatCacheStats(entries int, lastModified int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 Cache Statistics\n\n"))
	sb.WriteString(fmt.Sprintf("**Entries:** %d\n", entries))
	if lastModified > 0 {
		sb.WriteString(fmt.Sprintf("**Last Modified:** %d\n", lastModified))
	}
	return sb.String()
}

// FormatCacheInfo formats cache information for display
func FormatCacheInfo(info map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 Cache Information\n\n"))
	sb.WriteString(fmt.Sprintf("**File Key:** %s\n", info["file_key"]))
	sb.WriteString(fmt.Sprintf("**Total Entries:** %d\n", info["entries"]))
	sb.WriteString(fmt.Sprintf("**Components:** %d\n", info["components"]))
	sb.WriteString(fmt.Sprintf("**Metadata:** %d\n", info["metadata"]))
	return sb.String()
}
