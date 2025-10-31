package figma

import (
	"net/url"
	"regexp"
	"strings"
)

// IsFigmaURL checks if input contains figma.com
func IsFigmaURL(input string) bool {
	return strings.Contains(input, "figma.com")
}

// ParseFigmaURL extracts node-id from Figma URL
// Handles: node-id=123%3A456 or node-id=123:456
func ParseFigmaURL(urlStr string) (string, bool) {
	// Regex: node-id=([^&]+)
	// URL decode %3A → :
	re := regexp.MustCompile(`node-id=([^&]+)`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return "", false
	}
	nodeID, err := url.QueryUnescape(matches[1])
	if err != nil {
		return matches[1], true // Return raw if decode fails
	}
	return nodeID, true
}

// FilterComponents filters components by query (case-insensitive)
func FilterComponents(components []Component, query string) []Component {
	query = strings.ToLower(query)
	var results []Component
	for _, c := range components {
		if strings.Contains(strings.ToLower(c.Name), query) {
			results = append(results, c)
		}
	}
	return results
}

// GenerateComponentURL generates Figma URL for component
func GenerateComponentURL(fileKey, nodeID string) string {
	// Use string concatenation to avoid format string interpretation
	return "https://www.figma.com/file/" + fileKey + "/?node-id=" + nodeID
}

// TraverseTree recursively extracts components from file tree
func TraverseTree(node FigmaNode, fileKey string) []Component {
	var components []Component

	// Check if current node is a component
	if node.Type == "COMPONENT" || node.Type == "COMPONENT_SET" {
		components = append(components, Component{
			NodeID: node.ID,
			Name:   node.Name,
			Type:   node.Type,
			URL:    GenerateComponentURL(fileKey, node.ID),
		})
	}

	// Recurse into children
	for _, child := range node.Children {
		components = append(components, TraverseTree(child, fileKey)...)
	}

	return components
}
