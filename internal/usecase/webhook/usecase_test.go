package webhook_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Integral-07/prism/internal/domain/entity"
	"github.com/Integral-07/prism/internal/usecase/webhook"
)

type mockPRRepo struct {
	diff string
	err  error
}

func (m *mockPRRepo) GetDiff(_ context.Context, _ int64, _ string, _ int) (string, error) {
	return m.diff, m.err
}

type mockAnalyzerRepo struct {
	result webhook.AnalysisResult
	err    error
}

func (m *mockAnalyzerRepo) Analyze(_ context.Context, _, _ string) (webhook.AnalysisResult, error) {
	return m.result, m.err
}

type mockCheckRepo struct {
	called bool
	err    error
}

func (m *mockCheckRepo) PostResult(_ context.Context, _ int64, _ string, _ int, _ webhook.AnalysisResult) error {
	m.called = true
	return m.err
}

func TestInteractor_Execute(t *testing.T) {
	ctx := context.Background()
	input := webhook.Input{
		InstallationID: 123,
		RepoFullName:   "owner/repo",
		PRNumber:       1,
		Title:          "feat: add feature",
	}
	okResult := webhook.AnalysisResult{
		RiskLevel:        entity.RiskLevelLow,
		PriorityScore:    1,
		EstimatedMinutes: 5,
	}

	tests := []struct {
		name        string
		pr          *mockPRRepo
		analyzer    *mockAnalyzerRepo
		check       *mockCheckRepo
		wantErr     bool
		wantChecked bool
	}{
		{
			name:        "success",
			pr:          &mockPRRepo{diff: "diff content"},
			analyzer:    &mockAnalyzerRepo{result: okResult},
			check:       &mockCheckRepo{},
			wantErr:     false,
			wantChecked: true,
		},
		{
			name:        "pr repo error",
			pr:          &mockPRRepo{err: errors.New("github error")},
			analyzer:    &mockAnalyzerRepo{},
			check:       &mockCheckRepo{},
			wantErr:     true,
			wantChecked: false,
		},
		{
			name:        "analyzer error",
			pr:          &mockPRRepo{diff: "diff content"},
			analyzer:    &mockAnalyzerRepo{err: errors.New("claude error")},
			check:       &mockCheckRepo{},
			wantErr:     true,
			wantChecked: false,
		},
		{
			name:        "check repo error",
			pr:          &mockPRRepo{diff: "diff content"},
			analyzer:    &mockAnalyzerRepo{result: okResult},
			check:       &mockCheckRepo{err: errors.New("check error")},
			wantErr:     true,
			wantChecked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := webhook.NewInteractor(tt.pr, tt.analyzer, tt.check)
			err := uc.Execute(ctx, input)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if tt.check.called != tt.wantChecked {
				t.Errorf("wantChecked=%v, got %v", tt.wantChecked, tt.check.called)
			}
		})
	}
}
