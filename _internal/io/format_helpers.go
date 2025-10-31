package io

import (
	"fmt"
	"regexp"
	"strings"
)

// FormatWithEmoji adds contextual emoji based on content type
func FormatWithEmoji(text, contentType string) string {
	emoji := getEmojiForType(contentType)
	if emoji != "" {
		return fmt.Sprintf("%s %s", emoji, text)
	}
	return text
}

// Auto-detect and add emoji for common patterns
func AutoFormatEmoji(text string) string {
	// Detect URLs
	if strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://") {
		return FormatWithEmoji(text, "url")
	}
	// Detect ticket keys (PROJECT-123 pattern)
	if matched, _ := regexp.MatchString(`^[A-Z]+-\d+`, text); matched {
		return FormatWithEmoji(text, "ticket")
	}
	return text
}

func getEmojiForType(contentType string) string {
	emojiMap := map[string]string{
		"ticket":   "🎫",
		"url":      "🔗",
		"fetch":    "🌐",
		"cache":    "🔄",
		"success":  "✅",
		"error":    "❌",
		"warning":  "⚠️",
		"info":     "ℹ️",
		"file":     "📋",
		"commits":  "*",
		"branches": "⤴️",
		"user":     "👥",
		"activity": "👥",
		"ai":       "🤖",
	}
	return emojiMap[contentType]
}
