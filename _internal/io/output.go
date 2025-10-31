package io

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cli-go/_internal/ai"
	"cli-go/_internal/sys"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/glamour"
)

// Tools use DirectOutput() for all output handling

// ReadJSON reads JSON from stdin and unmarshals it into v
func ReadJSON(v interface{}) error {
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// WriteJSON marshals v to JSON and writes to stdout
func WriteJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// WriteErrorJSON writes an error as JSON to stdout
func WriteErrorJSON(err error) error {
	return WriteJSON(map[string]string{"error": err.Error()})
}

// DirectOutput handles output using direct flags (--clip, --file, --json)
// This replaces the old OutputRouter pattern with a simpler direct approach
func DirectOutput(data interface{}, clip bool, file string, asJSON bool) {
	var content string
	if asJSON {
		content = DataToString(data)
	} else {
		// For markdown/text content
		if str, ok := data.(string); ok {
			content = str
		} else {
			content = DataToString(data)
		}
	}

	// Write to destination(s) - support multiple
	if clip {
		ai.ExitIf(sys.CopyToClipboard(content), "clipboard error")
		fmt.Fprintf(os.Stderr, "Copied to clipboard\n")
	}
	if file != "" {
		ai.ExitIf(os.WriteFile(file, []byte(content), 0644), "file error")
		fmt.Fprintf(os.Stderr, "Written to %s\n", file)
	}
	// Always write to stdout if no explicit destination
	if !clip && file == "" {
		if asJSON {
			FormatTerminalOutputJSON(data)
		} else {
			FormatTerminalOutput(content)
		}
	}
}

// FormatTerminalOutput formats markdown using glamour and prints to terminal
func FormatTerminalOutput(content string) {
	// Create glamour renderer with dark theme
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		// Fallback to plain text if glamour fails
		fmt.Print(content)
		return
	}

	// Render markdown content
	output, err := r.Render(content)
	if err != nil {
		// Fallback to plain text if rendering fails
		fmt.Print(content)
		return
	}

	// Ensure proper newline
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}
	fmt.Print(output)
}

// FormatTerminalOutputWithResponseInfo formats markdown and adds response info
func FormatTerminalOutputWithResponseInfo(content string, responseInfo string) {
	// Create glamour renderer with dark theme
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		// Fallback to plain text if glamour fails
		fmt.Print(content)
		fmt.Fprintf(os.Stderr, "%s\n", responseInfo)
		return
	}

	// Render markdown content
	output, err := r.Render(content)
	if err != nil {
		// Fallback to plain text if rendering fails
		fmt.Print(content)
		fmt.Fprintf(os.Stderr, "%s\n", responseInfo)
		return
	}

	// Ensure proper newline
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}
	fmt.Print(output)

	// Add response info to stderr
	fmt.Fprintf(os.Stderr, "%s\n", responseInfo)
}

// FormatTerminalOutputJSON formats JSON using chroma and prints to terminal
func FormatTerminalOutputJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Use chroma for JSON syntax highlighting
	highlighted, err := highlightJSON(string(jsonData))
	if err != nil {
		// Fallback to raw JSON if chroma fails
		jsonStr := string(jsonData)
		if !strings.HasSuffix(jsonStr, "\n") {
			jsonStr += "\n"
		}
		fmt.Print(jsonStr)
		return nil
	}

	// Ensure proper newline
	if !strings.HasSuffix(highlighted, "\n") {
		highlighted += "\n"
	}
	fmt.Print(highlighted)
	return nil
}

// WriteToFile writes content to a file
func WriteToFile(content string, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// WriteToClipboard copies content to system clipboard (cross-platform)
func WriteToClipboard(content string) error {
	return sys.CopyToClipboard(content)
}

// DataToString converts interface{} to JSON string using chroma
func DataToString(data interface{}) string {
	if str, ok := data.(string); ok {
		return str
	}
	// Format JSON structure first
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	// Use chroma for syntax highlighting
	highlighted, err := highlightJSON(string(jsonData))
	if err != nil {
		// If chroma fails, return formatted JSON without highlighting
		return string(jsonData)
	}
	return highlighted
}

// highlightJSON uses chroma to highlight JSON syntax
func highlightJSON(jsonStr string) (string, error) {
	lexer := lexers.Get("json")
	if lexer == nil {
		return jsonStr, fmt.Errorf("json lexer not found")
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		return jsonStr, fmt.Errorf("terminal256 formatter not found")
	}

	style := styles.Get("monokai")
	if style == nil {
		return jsonStr, fmt.Errorf("monokai style not found")
	}

	iterator, err := lexer.Tokenise(nil, jsonStr)
	if err != nil {
		return jsonStr, err
	}

	var buf strings.Builder
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return jsonStr, err
	}

	return buf.String(), nil
}
