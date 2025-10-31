package ai

import (
	"net/http"
	"os"
)

// ExitIf handles the most common pattern: check error and exit if not nil
func ExitIf(err error, message string) {
	if err != nil {
		LogError("%s: %v", message, err)
		os.Exit(1)
	}
}

// HTTPIf handles the most common HTTP pattern: check error and send HTTP error if not nil
func HTTPIf(w http.ResponseWriter, err error, message string, statusCode int) {
	if err != nil {
		LogError("%s: %v", message, err)
		http.Error(w, message, statusCode)
	}
}
