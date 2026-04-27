package client

import (
	"fmt"
	"os"
	"strings"
)

// tokenFile файл хранения токена
const tokenFile = ".gophkeeper_token"

// SaveToken сохраняет JWT
func SaveToken(token string) error {
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// LoadToken загружает JWT
func LoadToken() (string, error) {
	b, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("not logged in")
	}
	return strings.TrimSpace(string(b)), nil
}
