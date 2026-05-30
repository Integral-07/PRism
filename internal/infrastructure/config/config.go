package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	GitHubAppID         int64
	GitHubPrivateKey    string
	GitHubWebhookSecret string
	GeminiAPIKey        string
	MockLLM             bool
	Port                string
}

func Load() (*Config, error) {
	appID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("GITHUB_APP_ID: %w", err)
	}

	privateKey := strings.ReplaceAll(os.Getenv("GITHUB_PRIVATE_KEY"), `\n`, "\n")
	if privateKey == "" {
		return nil, fmt.Errorf("GITHUB_PRIVATE_KEY is required")
	}

	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return nil, fmt.Errorf("GITHUB_WEBHOOK_SECRET is required")
	}

	mockLLM := os.Getenv("MOCK_LLM") == "true"

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" && !mockLLM {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		GitHubAppID:         appID,
		GitHubPrivateKey:    privateKey,
		GitHubWebhookSecret: webhookSecret,
		GeminiAPIKey:        geminiKey,
		MockLLM:             mockLLM,
		Port:                port,
	}, nil
}
