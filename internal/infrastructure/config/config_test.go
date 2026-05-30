package config_test

import (
	"testing"

	"github.com/Integral-07/prism/internal/infrastructure/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
		check   func(*testing.T, *config.Config)
	}{
		{
			name: "success with default port",
			env: map[string]string{
				"GITHUB_APP_ID":      "42",
				"GITHUB_PRIVATE_KEY": "my-key",
				"GEMINI_API_KEY":     "gemini-key",
			},
			check: func(t *testing.T, c *config.Config) {
				if c.GitHubAppID != 42 {
					t.Errorf("GitHubAppID: want 42, got %d", c.GitHubAppID)
				}
				if c.GitHubPrivateKey != "my-key" {
					t.Errorf("GitHubPrivateKey: want %q, got %q", "my-key", c.GitHubPrivateKey)
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
			env: map[string]string{
				"GITHUB_APP_ID":      "1",
				"GITHUB_PRIVATE_KEY": "key",
				"GEMINI_API_KEY":     "key",
				"PORT":               "9090",
			},
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
			name: "non-numeric GITHUB_APP_ID",
			env: map[string]string{
				"GITHUB_APP_ID": "abc",
			},
			wantErr: true,
		},
		{
			name: "missing GITHUB_PRIVATE_KEY",
			env: map[string]string{
				"GITHUB_APP_ID": "1",
			},
			wantErr: true,
		},
		{
			name: "missing GEMINI_API_KEY",
			env: map[string]string{
				"GITHUB_APP_ID":      "1",
				"GITHUB_PRIVATE_KEY": "key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range []string{"GITHUB_APP_ID", "GITHUB_PRIVATE_KEY", "GEMINI_API_KEY", "PORT"} {
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
