package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var groups = []string{"git", "ai", "tools"}

func isMainPackage(dir string) bool {
	ents, _ := os.ReadDir(dir)
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".go") {
			b, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err == nil && bytes.Contains(b, []byte("package main")) {
				return true
			}
		}
	}
	return false
}

// needsRebuild checks if a tool needs rebuilding based on file modification times
func needsRebuild(toolDir, outDir, toolName string) bool {
	outputPath := filepath.Join(outDir, toolName)

	// Check if output binary exists
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		return true // Binary doesn't exist, needs build
	}

	// Get all Go files that this tool depends on
	goFiles := getAllGoFiles(toolDir)

	// Check if any dependency file is newer than the binary
	for _, goFile := range goFiles {
		fileInfo, err := os.Stat(goFile)
		if err != nil {
			continue
		}

		if fileInfo.ModTime().After(outputInfo.ModTime()) {
			return true // Dependency file is newer than binary
		}
	}

	return false // Binary is up to date
}

// getAllGoFiles recursively finds all Go files that a tool depends on
func getAllGoFiles(toolDir string) []string {
	var goFiles []string

	// Get all Go files in the tool directory
	entries, err := os.ReadDir(toolDir)
	if err != nil {
		return goFiles
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			goFiles = append(goFiles, filepath.Join(toolDir, entry.Name()))
		}
	}

	// Get only the _internal packages that this tool actually imports
	internalFiles := getRelevantInternalFiles(toolDir)
	goFiles = append(goFiles, internalFiles...)

	// Also check go.mod and go.sum for dependency changes
	rootDir := filepath.Join(filepath.Dir(toolDir), "..")
	goMod := filepath.Join(rootDir, "go.mod")
	goSum := filepath.Join(rootDir, "go.sum")

	if _, err := os.Stat(goMod); err == nil {
		goFiles = append(goFiles, goMod)
	}
	if _, err := os.Stat(goSum); err == nil {
		goFiles = append(goFiles, goSum)
	}

	return goFiles
}

// getRelevantInternalFiles finds only the _internal files that a tool actually imports
func getRelevantInternalFiles(toolDir string) []string {
	var goFiles []string

	// Read the main.go file to find imports
	mainFile := filepath.Join(toolDir, "main.go")
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return goFiles
	}

	// Find all _internal imports
	lines := strings.Split(string(content), "\n")
	var internalPackages []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "cli-go/_internal/") {
			// Extract the package name from import
			parts := strings.Split(line, "cli-go/_internal/")
			if len(parts) > 1 {
				pkg := strings.Trim(parts[1], "\"")
				internalPackages = append(internalPackages, pkg)
			}
		}
	}

	// Get files for each imported _internal package
	internalDir := filepath.Join(filepath.Dir(toolDir), "..", "_internal")
	for _, pkg := range internalPackages {
		pkgDir := filepath.Join(internalDir, pkg)
		pkgFiles := getGoFilesRecursive(pkgDir)
		goFiles = append(goFiles, pkgFiles...)
		
		// Check for embedded filesystem assets in Go files
		embeddedAssets := getEmbeddedAssets(pkgDir)
		goFiles = append(goFiles, embeddedAssets...)
	}

	return goFiles
}

// getGoFilesRecursive recursively finds all Go files in a directory
func getGoFilesRecursive(dir string) []string {
	var goFiles []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return goFiles
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			// Recursively search subdirectories
			subFiles := getGoFilesRecursive(path)
			goFiles = append(goFiles, subFiles...)
		} else if strings.HasSuffix(entry.Name(), ".go") {
			goFiles = append(goFiles, path)
		}
	}

	return goFiles
}

// getEmbeddedAssets finds files referenced in //go:embed directives
func getEmbeddedAssets(pkgDir string) []string {
	var assets []string
	
	// Get all Go files in the package
	goFiles := getGoFilesRecursive(pkgDir)
	
	for _, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		if err != nil {
			continue
		}
		
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//go:embed") {
				// Extract the pattern from the embed directive
				pattern := strings.TrimSpace(strings.TrimPrefix(line, "//go:embed"))
				if pattern == "" {
					continue
				}
				
				// Convert the pattern to actual file paths
				embedFiles := expandEmbedPattern(filepath.Dir(goFile), pattern)
				assets = append(assets, embedFiles...)
				
				// Also check the next line for multi-line embed directives
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					if nextLine != "" && !strings.HasPrefix(nextLine, "//") && !strings.HasPrefix(nextLine, "var ") {
						// This might be a continuation of the embed pattern
						multiPattern := strings.TrimSpace(strings.TrimPrefix(nextLine, "//go:embed"))
						if multiPattern != "" {
							multiFiles := expandEmbedPattern(filepath.Dir(goFile), multiPattern)
							assets = append(assets, multiFiles...)
						}
					}
				}
			}
		}
	}
	
	return assets
}

// expandEmbedPattern expands a Go embed pattern to actual file paths
func expandEmbedPattern(baseDir, pattern string) []string {
	var files []string
	
	// Handle common embed patterns
	if strings.Contains(pattern, "*") {
		// Handle wildcard patterns like "web/dist/*"
		pattern = strings.TrimSpace(pattern)
		if strings.HasSuffix(pattern, "/*") {
			dir := strings.TrimSuffix(pattern, "/*")
			dirPath := filepath.Join(baseDir, dir)
			if entries, err := os.ReadDir(dirPath); err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						files = append(files, filepath.Join(dirPath, entry.Name()))
					} else {
						// Recursively include files in subdirectories
						subFiles := expandEmbedPattern(dirPath, entry.Name()+"/*")
						files = append(files, subFiles...)
					}
				}
			}
		}
	} else {
		// Handle specific file patterns
		filePath := filepath.Join(baseDir, pattern)
		if info, err := os.Stat(filePath); err == nil {
			if info.IsDir() {
				// If it's a directory, include all files in it
				if entries, err := os.ReadDir(filePath); err == nil {
					for _, entry := range entries {
						if !entry.IsDir() {
							files = append(files, filepath.Join(filePath, entry.Name()))
						} else {
							// Recursively include files in subdirectories
							subFiles := expandEmbedPattern(filePath, entry.Name()+"/*")
							files = append(files, subFiles...)
						}
					}
				}
			} else {
				// It's a file
				files = append(files, filePath)
			}
		}
	}
	
	return files
}

func main() {
	outDir := flag.String("o", "bin", "output directory")
	force := flag.Bool("force", false, "force rebuild all tools (ignore timestamps)")
	flag.Parse()
	_ = os.MkdirAll(*outDir, 0o755)

	root, _ := os.Getwd()
	buildErrors := 0

	for _, g := range groups {
		base := filepath.Join(root, g)
		entries, err := os.ReadDir(base)
		if err != nil {
			fmt.Printf("Warning: could not read %s directory: %v\n", g, err)
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			toolDir := filepath.Join(base, entry.Name())
			if !isMainPackage(toolDir) {
				continue
			}

			// Check if we should build this tool (always check unless --force)
			if !*force && !needsRebuild(toolDir, *outDir, entry.Name()) {
				fmt.Printf("skip  %-30s (up to date)\n", toolDir)
				continue
			}

			out := filepath.Join(*outDir, entry.Name())
			fmt.Printf("build %-30s -> %s\n", toolDir, out)

			cmd := exec.Command("go", "build", "-o", out, "./"+g+"/"+entry.Name())
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "build failed: %s: %v\n", entry.Name(), err)
				buildErrors++
			}
		}
	}

	if buildErrors > 0 {
		fmt.Fprintf(os.Stderr, "Build completed with %d errors\n", buildErrors)
		os.Exit(1)
	}

	fmt.Println("Build completed successfully")
}
