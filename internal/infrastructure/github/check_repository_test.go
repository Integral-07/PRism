package github

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Integral-07/prism/internal/domain/entity"
	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
	gh "github.com/google/go-github/v68/github"
)

func prHandler(sha string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"number": 1,
			"head":   map[string]string{"sha": sha},
		})
	}
}

func checkRunHandler(t *testing.T, wantConclusion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()
		var body gh.CreateCheckRunOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if body.Name != "PRism" {
			t.Errorf("check run name: want PRism, got %s", body.Name)
		}
		if body.HeadSHA != "abc123" {
			t.Errorf("HeadSHA: want abc123, got %s", body.HeadSHA)
		}
		if wantConclusion != "" {
			if body.GetConclusion() != wantConclusion {
				t.Errorf("conclusion: want %q, got %q", wantConclusion, body.GetConclusion())
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": 1, "name": "PRism"})
	}
}

func TestCheckRepository_PostResult(t *testing.T) {
	ctx := context.Background()
	defaultCfg := entity.DefaultPrismConfig()

	tests := []struct {
		name           string
		clientFor      func(int64) (*gh.Client, error)
		prHandler      http.HandlerFunc
		checkHandler   http.HandlerFunc
		result         analyzepr.Output
		wantConclusion string
		wantErr        bool
	}{
		{
			name:           "success low risk → success conclusion",
			result:         analyzepr.Output{RiskLevel: entity.RiskLevelLow, PriorityScore: 2, EstimatedMinutes: 10},
			wantConclusion: "success",
		},
		{
			name:           "success medium risk → success conclusion",
			result:         analyzepr.Output{RiskLevel: entity.RiskLevelMedium, PriorityScore: 3, EstimatedMinutes: 30},
			wantConclusion: "success",
		},
		{
			name:           "success high risk → neutral conclusion",
			result:         analyzepr.Output{RiskLevel: entity.RiskLevelHigh, PriorityScore: 5, EstimatedMinutes: 60},
			wantConclusion: "neutral",
		},
		{
			name:      "client error",
			clientFor: errorClient(errors.New("auth error")),
			wantErr:   true,
		},
		{
			name: "get PR error",
			prHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "not found", http.StatusNotFound)
			},
			wantErr: true,
		},
		{
			name: "create check run error",
			checkHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "forbidden", http.StatusForbidden)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf func(int64) (*gh.Client, error)
			if tt.clientFor != nil {
				cf = tt.clientFor
			} else {
				mux := http.NewServeMux()

				ph := tt.prHandler
				if ph == nil {
					ph = prHandler("abc123")
				}
				mux.HandleFunc("/repos/owner/repo/pulls/1", ph)

				ch := tt.checkHandler
				if ch == nil {
					ch = checkRunHandler(t, tt.wantConclusion)
				}
				mux.HandleFunc("/repos/owner/repo/check-runs", ch)

				cf = fixedClient(testGHClient(t, mux))
			}

			repo := &CheckRepository{clientFor: cf}
			err := repo.PostResult(ctx, 123, "owner/repo", 1, tt.result, defaultCfg)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestBuildSummary(t *testing.T) {
	defaultCfg := entity.DefaultPrismConfig()

	tests := []struct {
		name   string
		output analyzepr.Output
		cfg    entity.PrismConfig
		checks []string
	}{
		{
			name: "summary and review focus shown",
			output: analyzepr.Output{
				RiskLevel:        entity.RiskLevelLow,
				Summary:          "軽微な変更です",
				ReviewFocus:      []string{"影響範囲を確認"},
				EstimatedMinutes: 10,
			},
			cfg:    defaultCfg,
			checks: []string{"軽微な変更です", "影響範囲を確認"},
		},
		{
			name: "file priority list shown",
			output: analyzepr.Output{
				RiskLevel: entity.RiskLevelHigh,
				Files: []analyzepr.FileRisk{
					{Path: "auth/jwt.go", Risk: entity.RiskLevelHigh, Category: "logic", Note: "要確認"},
				},
				EstimatedMinutes: 60,
			},
			cfg:    defaultCfg,
			checks: []string{"auth/jwt.go", "logic", "要確認"},
		},
		{
			name: "custom output shown",
			output: analyzepr.Output{
				RiskLevel:    entity.RiskLevelLow,
				CustomOutput: "決済ロジックへの影響なし",
			},
			cfg:    defaultCfg,
			checks: []string{"決済ロジックへの影響なし"},
		},
		{
			name: "review focus hidden when disabled",
			output: analyzepr.Output{
				RiskLevel:   entity.RiskLevelLow,
				ReviewFocus: []string{"注意点"},
			},
			cfg: entity.PrismConfig{
				Output: entity.OutputConfig{
					Support: entity.SupportConfig{ReviewFocus: false},
				},
			},
			checks: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSummary(tt.output, tt.cfg)
			for _, s := range tt.checks {
				if !strings.Contains(got, s) {
					t.Errorf("summary does not contain %q\ngot: %s", s, got)
				}
			}
			if tt.name == "review focus hidden when disabled" && strings.Contains(got, "注意点") {
				t.Errorf("summary should not contain hidden review focus")
			}
		})
	}
}
