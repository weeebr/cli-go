package figma

import (
	"fmt"
	"net/url"
	"strings"
)

// ExtractFileIDFromURL extracts the file ID from a Figma URL
func ExtractFileIDFromURL(figmaURL string) (string, error) {
	parsedURL, err := url.Parse(figmaURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Handle different Figma URL formats
	if strings.Contains(parsedURL.Host, "figma.com") {
		pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		if len(pathParts) >= 2 && pathParts[0] == "file" {
			return pathParts[1], nil
		}
	}

	return "", fmt.Errorf("could not extract file ID from URL")
}

// ExtractNodeIDFromURL extracts the node ID from a Figma URL
func ExtractNodeIDFromURL(figmaURL string) (string, error) {
	parsedURL, err := url.Parse(figmaURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Look for node ID in URL fragment or query parameters
	if parsedURL.Fragment != "" {
		return parsedURL.Fragment, nil
	}

	// Check query parameters
	nodeID := parsedURL.Query().Get("node-id")
	if nodeID != "" {
		return nodeID, nil
	}

	return "", fmt.Errorf("could not extract node ID from URL")
}

// BuildFigmaURL constructs a Figma URL from file ID and node ID
func BuildFigmaURL(fileID, nodeID string) string {
	if nodeID != "" {
		return fmt.Sprintf("https://www.figma.com/file/%s?node-id=%s", fileID, nodeID)
	}
	return fmt.Sprintf("https://www.figma.com/file/%s", fileID)
}

// FormatComponentMetadata formats component metadata for display
func FormatComponentMetadata(metadata *EnhancedMetadata) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("## Component: %s\n", metadata.Component.Name))
	result.WriteString(fmt.Sprintf("**ID:** %s\n", metadata.Component.ID))
	result.WriteString(fmt.Sprintf("**Type:** %s\n", metadata.Component.Type))

	if metadata.Description != "" {
		result.WriteString(fmt.Sprintf("**Description:** %s\n", metadata.Description))
	}

	result.WriteString(fmt.Sprintf("**Properties:** %d\n", metadata.PropertyCount))
	result.WriteString(fmt.Sprintf("**Interactions:** %t\n", metadata.HasInteractions))
	result.WriteString(fmt.Sprintf("**Overrides:** %d\n", metadata.Overrides))

	if metadata.LayoutMode != "" {
		result.WriteString(fmt.Sprintf("**Layout Mode:** %s\n", metadata.LayoutMode))
	}

	if metadata.CornerRadius > 0 {
		result.WriteString(fmt.Sprintf("**Corner Radius:** %.2f\n", metadata.CornerRadius))
	}

	if metadata.Opacity > 0 {
		result.WriteString(fmt.Sprintf("**Opacity:** %.2f\n", metadata.Opacity))
	}

	if metadata.BlendMode != "" {
		result.WriteString(fmt.Sprintf("**Blend Mode:** %s\n", metadata.BlendMode))
	}

	result.WriteString(fmt.Sprintf("**Children Count:** %d\n", metadata.ChildrenCount))
	result.WriteString(fmt.Sprintf("**Visible:** %t\n", metadata.Visible))
	result.WriteString(fmt.Sprintf("**Locked:** %t\n", metadata.Locked))
	result.WriteString(fmt.Sprintf("**Is Mask:** %t\n", metadata.IsMask))

	return result.String()
}

// FormatProperties formats component properties for display
func FormatProperties(properties map[string]interface{}) string {
	if len(properties) == 0 {
		return "No properties defined"
	}

	var result strings.Builder
	result.WriteString("## Properties\n")

	for key, value := range properties {
		result.WriteString(fmt.Sprintf("- **%s:** %v\n", key, value))
	}

	return result.String()
}

// FormatInteractions formats component interactions for display
func FormatInteractions(interactions []interface{}) string {
	if len(interactions) == 0 {
		return "No interactions defined"
	}

	var result strings.Builder
	result.WriteString("## Interactions\n")

	for i, interaction := range interactions {
		result.WriteString(fmt.Sprintf("%d. %v\n", i+1, interaction))
	}

	return result.String()
}

// FormatStyleReferences formats style references for display
func FormatStyleReferences(styleRefs map[string]string) string {
	if len(styleRefs) == 0 {
		return "No style references"
	}

	var result strings.Builder
	result.WriteString("## Style References\n")

	for key, value := range styleRefs {
		result.WriteString(fmt.Sprintf("- **%s:** %s\n", key, value))
	}

	return result.String()
}

// FormatConstraints formats component constraints for display
func FormatConstraints(constraints map[string]interface{}) string {
	if len(constraints) == 0 {
		return "No constraints defined"
	}

	var result strings.Builder
	result.WriteString("## Constraints\n")

	for key, value := range constraints {
		result.WriteString(fmt.Sprintf("- **%s:** %v\n", key, value))
	}

	return result.String()
}

// FormatAutoLayout formats auto layout properties for display
func FormatAutoLayout(autoLayout map[string]interface{}) string {
	if len(autoLayout) == 0 {
		return "No auto layout properties"
	}

	var result strings.Builder
	result.WriteString("## Auto Layout\n")

	for key, value := range autoLayout {
		result.WriteString(fmt.Sprintf("- **%s:** %v\n", key, value))
	}

	return result.String()
}

// FormatSize formats component size information for display
func FormatSize(size map[string]interface{}) string {
	if len(size) == 0 {
		return "No size information"
	}

	var result strings.Builder
	result.WriteString("## Size\n")

	for key, value := range size {
		result.WriteString(fmt.Sprintf("- **%s:** %v\n", key, value))
	}

	return result.String()
}

// FormatPadding formats padding information for display
func FormatPadding(padding map[string]float64) string {
	if len(padding) == 0 {
		return "No padding defined"
	}

	var result strings.Builder
	result.WriteString("## Padding\n")

	for key, value := range padding {
		result.WriteString(fmt.Sprintf("- **%s:** %.2f\n", key, value))
	}

	return result.String()
}

// FormatEffects formats component effects for display
func FormatEffects(effects []interface{}) string {
	if len(effects) == 0 {
		return "No effects applied"
	}

	var result strings.Builder
	result.WriteString("## Effects\n")

	for i, effect := range effects {
		result.WriteString(fmt.Sprintf("%d. %v\n", i+1, effect))
	}

	return result.String()
}

// FormatFills formats component fills for display
func FormatFills(fills []interface{}) string {
	if len(fills) == 0 {
		return "No fills defined"
	}

	var result strings.Builder
	result.WriteString("## Fills\n")

	for i, fill := range fills {
		result.WriteString(fmt.Sprintf("%d. %v\n", i+1, fill))
	}

	return result.String()
}

// FormatStrokes formats component strokes for display
func FormatStrokes(strokes []interface{}) string {
	if len(strokes) == 0 {
		return "No strokes defined"
	}

	var result strings.Builder
	result.WriteString("## Strokes\n")

	for i, stroke := range strokes {
		result.WriteString(fmt.Sprintf("%d. %v\n", i+1, stroke))
	}

	return result.String()
}

// FormatTextProperties formats text properties for display
func FormatTextProperties(textProps map[string]interface{}) string {
	if len(textProps) == 0 {
		return "No text properties"
	}

	var result strings.Builder
	result.WriteString("## Text Properties\n")

	for key, value := range textProps {
		result.WriteString(fmt.Sprintf("- **%s:** %v\n", key, value))
	}

	return result.String()
}

// ExtractFileAndNodeFromURL extracts both file key and node ID from a Figma URL
func ExtractFileAndNodeFromURL(url string) (string, string, error) {
	fileKey, err := ExtractFileIDFromURL(url)
	if err != nil {
		return "", "", err
	}

	nodeID, err := ExtractNodeIDFromURL(url)
	if err != nil {
		return "", "", err
	}

	return fileKey, nodeID, nil
}
