package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func main() {
	fmt.Println("Go files found:")
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
