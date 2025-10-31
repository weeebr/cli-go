# Go CLI Tools

A collection of command-line tools built with Go.

<img width="1021" height="497" alt="screenshot" src="https://github.com/user-attachments/assets/d150b867-9919-4c38-9046-01f64495d8c5" />


## Structure

- `/<group>/<tool>/main.go` → builds to binary name
- Groups: `git/`, `ai/`, `tools/`, `_internal/`
- Internal libs in `_internal/` (no package main)

## Development Standards

- ≤300 LOC per tool
- No globals - use dependency injection
- No fixtures - validate via terminal commands only
- No test files - test functionality via terminal execution
- Never use timeout, always sleep in terminal executions
- Output: Human-friendly by default, JSON with `--json` flag
- Exit non-zero on error

## Output System

- **Default**: Human-friendly output with emojis
- **JSON**: Use `--json` flag for automation
- **Parsing**: Use `flags.ReorderAndParse()` instead of `flag.Parse()`

### Pattern

```go
import (
    "flag"
    "cli-go/_internal/flags"
    "cli-go/_internal/io"
)

func main() {
    var (
        output = flag.String("output", "stdout", "Output destination")
        path   = flag.String("path", "", "File path (only with --output=file)")
        json   = flag.Bool("json", false, "Output in JSON format")
    )
    flags.ReorderAndParse() // CRITICAL: Use this instead of flag.Parse()

    result := map[string]interface{}{"message": "Success", "success": true}

    if *json {
        config := io.OutputConfig{Mode: io.OutputMode(*output), Path: *path}
        io.NewOutputRouter(config).Write(result)
    } else {
        fmt.Printf("✅ %s\n", result["message"])
    }
}
```
