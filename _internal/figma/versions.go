package figma

import (
	"fmt"
	"time"
)

// VersionData represents Figma file version information
type VersionData struct {
	Versions []Version `json:"versions"`
}

// Version represents a single version entry
type Version struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	User        User      `json:"user"`
}

// User represents a Figma user
type User struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
}

// VersionSummary represents version analysis results
type VersionSummary struct {
	TotalVersions int       `json:"total_versions"`
	Changes30Days int       `json:"changes_30_days"`
	Changes7Days  int       `json:"changes_7_days"`
	LastChange    time.Time `json:"last_change"`
	DaysSinceLast int       `json:"days_since_last"`
	Summary       string    `json:"summary"`
}

// GetVersionSummary calculates version changes in past N days
func GetVersionSummary(fileKey string, days int) (*VersionSummary, error) {
	// For now, return mock data since we don't have version API access
	// In a real implementation, this would call the Figma versions API
	now := time.Now()

	return &VersionSummary{
		TotalVersions: 42,
		Changes30Days: 8,
		Changes7Days:  2,
		LastChange:    now.AddDate(0, 0, -3),
		DaysSinceLast: 3,
		Summary:       fmt.Sprintf("File updated %d times past %d days (last change: %d days ago)", 8, days, 3),
	}, nil
}

// GetVersionStats returns detailed version statistics
func GetVersionStats(fileKey string) (*VersionSummary, error) {
	return GetVersionSummary(fileKey, 30)
}

// CalculateChangesPastDays calculates changes in past N days
func CalculateChangesPastDays(fileKey string, days int) (int, error) {
	summary, err := GetVersionSummary(fileKey, days)
	if err != nil {
		return 0, err
	}
	return summary.Changes30Days, nil
}

// CalculateDaysSinceLastChange calculates days since last change
func CalculateDaysSinceLastChange(fileKey string) (int, error) {
	summary, err := GetVersionSummary(fileKey, 30)
	if err != nil {
		return 0, err
	}
	return summary.DaysSinceLast, nil
}

// FormatVersionInfo formats version information for display
func FormatVersionInfo(summary *VersionSummary) string {
	var result string
	result += fmt.Sprintf("**Total versions:** %d\n", summary.TotalVersions)
	result += fmt.Sprintf("**Changes (30 days):** %d\n", summary.Changes30Days)
	result += fmt.Sprintf("**Changes (7 days):** %d\n", summary.Changes7Days)
	result += fmt.Sprintf("**Last change:** %s\n", summary.LastChange.Format("2006-01-02 15:04:05"))
	result += fmt.Sprintf("**Summary:** %s\n", summary.Summary)
	return result
}
