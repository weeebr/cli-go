package custom

import (
	"fmt"
	"regexp"
)

// NormalizeTicketID converts numeric shorthand to full issue key
// This consolidates logic from both Jira and GitHub tools
func NormalizeTicketID(ticketID, defaultProjectKey string) string {
	if ticketID == "" {
		return ""
	}

	// Check if already in PROJECT-123 format
	matched, _ := regexp.MatchString(`^[A-Z]{2,}-\d+$`, ticketID)
	if matched {
		return ticketID
	}

	// Check if it's just a number
	matched, _ = regexp.MatchString(`^\d+$`, ticketID)
	if matched {
		return defaultProjectKey + "-" + ticketID
	}

	// Invalid format - return as-is (caller should validate)
	return ticketID
}

// ValidateTicketID validates and normalizes a ticket ID
func ValidateTicketID(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("ticket ID cannot be empty")
	}

	// Check if it's a valid format
	validPattern := regexp.MustCompile(`^[A-Z]{2,}-\d+$|^\d+$`)
	if !validPattern.MatchString(input) {
		return "", fmt.Errorf("invalid ticket ID format: %s", input)
	}

	return input, nil
}

// ExtractTicketsFromMessage extracts ticket IDs from a commit message
func ExtractTicketsFromMessage(message string) []string {
	// Universal ticket detection regex (PNT-123, RDP-456, CASH-789, ORB-321, etc.)
	re := regexp.MustCompile(`\b([A-Z]{2,}-\d+)\b`)
	matches := re.FindAllString(message, -1)

	// Remove duplicates
	seen := make(map[string]bool)
	var tickets []string
	for _, match := range matches {
		if !seen[match] {
			seen[match] = true
			tickets = append(tickets, match)
		}
	}

	return tickets
}

// ParseTicketPattern returns the compiled regex for ticket matching
func ParseTicketPattern() *regexp.Regexp {
	return regexp.MustCompile(`\b([A-Z]{2,}-\d+)\b`)
}

// IsValidTicketFormat checks if a string matches the ticket format
func IsValidTicketFormat(ticketID string) bool {
	pattern := ParseTicketPattern()
	return pattern.MatchString(ticketID)
}

// GetTicketFromBranch extracts ticket ID from a branch name
func GetTicketFromBranch(branchName string) string {
	// Common branch patterns: feature/PNT-123, PNT-123, bugfix/PNT-123
	pattern := regexp.MustCompile(`(?:feature/|bugfix/|hotfix/)?([A-Z]{2,}-\d+)`)
	matches := pattern.FindStringSubmatch(branchName)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
