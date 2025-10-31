// ALL HERE IS JUST READ-ONLY ACCESS. NO EDITS, NO ADD_KEY, NO SETUP.

package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// getCredentialsFile returns the credentials file path from config
func getCredentialsFile() string {
	cfg, _ := LoadConfig()

	return cfg.Credentials.File
}

type EncryptedStore struct {
	filePath string
}

type Credentials struct {
	OpenAI     string `json:"openai"`
	Anthropic  string `json:"anthropic"`
	Google     string `json:"google"`
	Groq       string `json:"groq"`
	Perplexity string `json:"perplexity"`
	Figma      string `json:"figma"`
	Jira       string `json:"jira"`
}

func NewEncryptedStore() (*EncryptedStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	filePath := filepath.Join(homeDir, getCredentialsFile())
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	return &EncryptedStore{filePath: filePath}, nil
}

func getMachineKey() []byte {
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	data := fmt.Sprintf("%s:%s:cli-go", hostname, username)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// SaveCredentials - ONLY used by setup tool
func (es *EncryptedStore) SaveCredentials(creds *Credentials) error {
	jsonData, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %v", err)
	}

	key := getMachineKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)
	return os.WriteFile(es.filePath, ciphertext, 0600)
}

// LoadCredentials - internal helper
func (es *EncryptedStore) LoadCredentials() (*Credentials, error) {
	data, err := os.ReadFile(es.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("credentials not found")
		}
		return nil, fmt.Errorf("failed to read credentials: %v", err)
	}

	key := getMachineKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("invalid credentials file")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt")
	}

	var creds Credentials
	if err := json.Unmarshal(plaintext, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	return &creds, nil
}

func (es *EncryptedStore) GetKey(service string) (string, error) {
	creds, err := es.LoadCredentials()
	if err != nil {
		return "", err
	}

	switch service {
	case "openai":
		if creds.OpenAI == "" {
			return "", fmt.Errorf("openai key not configured")
		}
		return creds.OpenAI, nil
	case "anthropic":
		if creds.Anthropic == "" {
			return "", fmt.Errorf("anthropic key not configured")
		}
		return creds.Anthropic, nil
	case "google":
		if creds.Google == "" {
			return "", fmt.Errorf("google key not configured")
		}
		return creds.Google, nil
	case "groq":
		if creds.Groq == "" {
			return "", fmt.Errorf("groq key not configured")
		}
		return creds.Groq, nil
	case "perplexity":
		if creds.Perplexity == "" {
			return "", fmt.Errorf("perplexity key not configured")
		}
		return creds.Perplexity, nil
	case "jira":
		if creds.Jira == "" {
			return "", fmt.Errorf("jira key not configured")
		}
		return creds.Jira, nil
	case "figma":
		if creds.Figma == "" {
			return "", fmt.Errorf("figma key not configured")
		}
		return creds.Figma, nil
	default:
		return "", fmt.Errorf("unknown service: %s", service)
	}
}

func (es *EncryptedStore) Exists() bool {
	_, err := os.Stat(es.filePath)
	return err == nil
}

// GetKey is a convenience function to get a key without creating a store
func GetKey(service string) (string, error) {
	store, err := NewEncryptedStore()
	if err != nil {
		return "", err
	}
	return store.GetKey(service)
}
