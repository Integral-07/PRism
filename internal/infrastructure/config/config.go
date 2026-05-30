package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GitHubAppID      int64
	GitHubPrivateKey string
	GeminiAPIKey     string
	Port             string
}

func Load() (*Config, error) {
	appID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("GITHUB_APP_ID: %w", err)
	}

	privateKey := os.Getenv("GITHUB_PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("GITHUB_PRIVATE_KEY is required")
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		GitHubAppID:      appID,
		GitHubPrivateKey: privateKey,
		GeminiAPIKey:     geminiKey,
		Port:             port,
	}, nil
}
