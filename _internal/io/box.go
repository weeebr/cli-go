package io

import (
	"fmt"
	"strings"
)

// FormatBoxed creates a bordered box around text
func FormatBoxed(title string) string {
	boxWidth := len(title) + 4

	var output strings.Builder
	output.WriteString("╭" + strings.Repeat("─", boxWidth-2) + "╮\n")
	output.WriteString("│" + strings.Repeat(" ", boxWidth-2) + "│\n")
	output.WriteString(fmt.Sprintf("│ %s │\n", title))
	output.WriteString("│" + strings.Repeat(" ", boxWidth-2) + "│\n")
	output.WriteString("╰" + strings.Repeat("─", boxWidth-2) + "╯")

	return output.String()
}
