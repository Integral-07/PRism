package config_test

import (
	"testing"

	"github.com/Integral-07/prism/internal/infrastructure/config"
)

func TestLoad(t *testing.T) {
	allKeys := []string{"GITHUB_APP_ID", "GITHUB_PRIVATE_KEY", "GITHUB_WEBHOOK_SECRET", "GEMINI_API_KEY", "PORT"}

	base := map[string]string{
		"GITHUB_APP_ID":          "42",
		"GITHUB_PRIVATE_KEY":     "my-key",
		"GITHUB_WEBHOOK_SECRET":  "webhook-secret",
		"GEMINI_API_KEY":         "gemini-key",
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
		check   func(*testing.T, *config.Config)
	}{
		{
			name: "success with default port",
			env:  base,
			check: func(t *testing.T, c *config.Config) {
				if c.GitHubAppID != 42 {
					t.Errorf("GitHubAppID: want 42, got %d", c.GitHubAppID)
				}
				if c.GitHubPrivateKey != "my-key" {
					t.Errorf("GitHubPrivateKey: want %q, got %q", "my-key", c.GitHubPrivateKey)
				}
				if c.GitHubWebhookSecret != "webhook-secret" {
					t.Errorf("GitHubWebhookSecret: want %q, got %q", "webhook-secret", c.GitHubWebhookSecret)
				}
				if c.GeminiAPIKey != "gemini-key" {
					t.Errorf("GeminiAPIKey: want %q, got %q", "gemini-key", c.GeminiAPIKey)
				}
				if c.Port != "8080" {
					t.Errorf("Port: want 8080, got %s", c.Port)
				}
			},
		},
		{
			name: "custom port",
			env:  merge(base, map[string]string{"PORT": "9090"}),
			check: func(t *testing.T, c *config.Config) {
				if c.Port != "9090" {
					t.Errorf("Port: want 9090, got %s", c.Port)
				}
			},
		},
		{
			name:    "missing GITHUB_APP_ID",
			env:     map[string]string{},
			wantErr: true,
		},
		{
			name:    "non-numeric GITHUB_APP_ID",
			env:     map[string]string{"GITHUB_APP_ID": "abc"},
			wantErr: true,
		},
		{
			name:    "missing GITHUB_PRIVATE_KEY",
			env:     map[string]string{"GITHUB_APP_ID": "1"},
			wantErr: true,
		},
		{
			name:    "missing GITHUB_WEBHOOK_SECRET",
			env:     map[string]string{"GITHUB_APP_ID": "1", "GITHUB_PRIVATE_KEY": "key"},
			wantErr: true,
		},
		{
			name:    "missing GEMINI_API_KEY",
			env:     map[string]string{"GITHUB_APP_ID": "1", "GITHUB_PRIVATE_KEY": "key", "GITHUB_WEBHOOK_SECRET": "secret"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range allKeys {
				t.Setenv(key, "")
			}
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func merge(base, override map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}
