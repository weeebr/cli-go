package jira

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"cli-go/_internal/io"
)

// ADFNode represents a node in Atlassian Document Format
type ADFNode struct {
	Type    string                 `json:"type"`
	Content []ADFNode              `json:"content,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Marks   []ADFMark              `json:"marks,omitempty"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
}

// ADFMark represents formatting marks in ADF
type ADFMark struct {
	Type string `json:"type"`
}

// ConvertADFToMarkdown converts ADF content to markdown using external library
func ConvertADFToMarkdown(adfContent interface{}) (string, error) {
	// Handle nil or empty content
	if adfContent == nil {
		return "No description", nil
	}

	// If it's already a string, return as-is
	if str, ok := adfContent.(string); ok {
		if str == "" || str == "null" {
			return "No description", nil
		}
		return str, nil
	}

	// Convert to JSON for processing
	jsonData, err := json.Marshal(adfContent)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ADF content: %v", err)
	}

	// Try to use adf2md library if available
	markdown, err := convertWithADF2MD(string(jsonData))
	if err == nil {
		return markdown, nil
	}

	// Fallback to simple conversion
	return convertADFSimple(adfContent)
}

// convertWithADF2MD uses external adf2md library
func convertWithADF2MD(jsonData string) (string, error) {
	// Check if adf2md is available
	if _, err := exec.LookPath("adf2md"); err != nil {
		return "", fmt.Errorf("adf2md not found")
	}

	cmd := exec.Command("adf2md")
	cmd.Stdin = strings.NewReader(jsonData)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("adf2md conversion failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// convertADFSimple provides a basic ADF to markdown conversion
func convertADFSimple(adfContent interface{}) (string, error) {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Panic recovered, will return error below
		}
	}()

	// Try to parse as ADF structure
	jsonData, err := json.Marshal(adfContent)
	if err != nil {
		return "", err
	}

	var node ADFNode
	if err := json.Unmarshal(jsonData, &node); err != nil {
		// If it's not ADF, return as string
		return string(jsonData), nil
	}

	markdown := convertNodeToMarkdown(node)

	return markdown, nil
}

// convertNodeToMarkdown converts an ADF node to markdown
func convertNodeToMarkdown(node ADFNode) string {
	// Add recursion limit to prevent infinite loops
	return convertNodeToMarkdownWithLimit(node, 0, 10, 0)
}

// convertNodeToMarkdownWithLimit converts an ADF node to markdown with recursion limit
func convertNodeToMarkdownWithLimit(node ADFNode, depth, maxDepth, listDepth int) string {
	// Prevent infinite recursion
	if depth > maxDepth {
		return "[...]"
	}
	switch node.Type {
	case "text":
		text := node.Text
		// Apply marks
		for _, mark := range node.Marks {
			switch mark.Type {
			case "strong":
				text = "**" + text + "**"
			case "em":
				text = "*" + text + "*"
			case "code":
				text = "`" + text + "`"
			}
		}
		return text

	case "paragraph":
		var content strings.Builder
		for _, child := range node.Content {
			content.WriteString(convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth))
		}
		// Clean up adjacent bold formatting
		text := content.String()
		// Fix double bold markers
		text = strings.ReplaceAll(text, "****", "")
		// Fix triple bold markers
		text = strings.ReplaceAll(text, "***", "*")

		// Only add newline if content is not empty
		cleanContent := strings.TrimSpace(text)
		if cleanContent == "" {
			return ""
		}

		// Check if this is a header item (no bullet, just text, and likely a component name)
		if !strings.Contains(cleanContent, "•") && !strings.Contains(cleanContent, "✅") && !strings.Contains(cleanContent, "❓") && !strings.Contains(cleanContent, "🔗") && !strings.Contains(cleanContent, "https://") {
			// This is likely a header item, don't add bullet
			return cleanContent + "\n"
		}

		return cleanContent + "\n"

	case "heading":
		var content strings.Builder
		for _, child := range node.Content {
			content.WriteString(convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth))
		}
		level := 1
		if node.Attrs != nil {
			if l, ok := node.Attrs["level"].(float64); ok {
				level = int(l)
			}
		}
		return strings.Repeat("#", level) + " " + content.String() + "\n"

	case "codeBlock":
		var content strings.Builder
		for _, child := range node.Content {
			content.WriteString(convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth))
		}
		language := ""
		if node.Attrs != nil {
			if lang, ok := node.Attrs["language"].(string); ok {
				language = lang
			}
		}
		return "```" + language + "\n" + content.String() + "\n```\n"

	case "bulletList":
		var content strings.Builder
		for i, child := range node.Content {
			item := convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth+1)
			content.WriteString(item)
			// Only add newline between items, not after the last one
			if i < len(node.Content)-1 {
				content.WriteString("\n")
			}
		}
		return content.String()

	case "listItem":
		var content strings.Builder
		for _, child := range node.Content {
			content.WriteString(convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth))
		}
		// Remove trailing newlines from content and add proper list formatting
		cleanContent := strings.TrimRight(content.String(), "\n")
		// Add indentation based on list depth
		indent := strings.Repeat("  ", listDepth)
		return indent + "• " + cleanContent

	case "orderedList":
		var content strings.Builder
		for i, child := range node.Content {
			item := convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth+1)
			cleanItem := strings.TrimRight(item, "\n")
			content.WriteString(fmt.Sprintf("%d. %s", i+1, cleanItem))
			// Only add newline between items, not after the last one
			if i < len(node.Content)-1 {
				content.WriteString("\n")
			}
		}
		return content.String()

	case "hardBreak":
		return "\n"

	case "mention":
		if node.Attrs != nil {
			if text, ok := node.Attrs["text"].(string); ok {
				return text
			}
			if displayName, ok := node.Attrs["displayName"].(string); ok {
				return displayName
			}
		}
		return "@unknown"

	case "link":
		if node.Attrs != nil {
			if href, ok := node.Attrs["href"].(string); ok {
				return io.AutoFormatEmoji(href)
			}
		}
		return io.FormatWithEmoji("unknown", "url")

	case "image":
		alt := "Image"
		src := ""
		if node.Attrs != nil {
			if altText, ok := node.Attrs["alt"].(string); ok {
				alt = altText
			}
			if srcURL, ok := node.Attrs["src"].(string); ok {
				src = srcURL
			}
		}
		if src != "" {
			return fmt.Sprintf("![%s](%s)", alt, src)
		}
		return "@image"

	case "emoji":
		if node.Attrs != nil {
			if text, ok := node.Attrs["text"].(string); ok {
				return text
			}
			if shortName, ok := node.Attrs["shortName"].(string); ok {
				return ":" + shortName + ":"
			}
		}
		return ":unknown:"

	case "inlineCard":
		if node.Attrs != nil {
			if url, ok := node.Attrs["url"].(string); ok {
				return io.AutoFormatEmoji(url)
			}
		}
		return io.FormatWithEmoji("unknown", "url")

	default:
		// For unknown types, try to process content
		var content strings.Builder
		for _, child := range node.Content {
			content.WriteString(convertNodeToMarkdownWithLimit(child, depth+1, maxDepth, listDepth))
		}
		return content.String()
	}
}
