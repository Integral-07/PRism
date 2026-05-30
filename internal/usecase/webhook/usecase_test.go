package webhook_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Integral-07/prism/internal/domain/entity"
	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
	"github.com/Integral-07/prism/internal/usecase/webhook"
)

type mockPRRepo struct {
	diff string
	err  error
}

func (m *mockPRRepo) GetDiff(_ context.Context, _ int64, _ string, _ int) (string, error) {
	return m.diff, m.err
}

type mockAnalyzerUC struct {
	result analyzepr.Output
	err    error
}

func (m *mockAnalyzerUC) Execute(_ context.Context, _ analyzepr.Input) (analyzepr.Output, error) {
	return m.result, m.err
}

type mockCheckRepo struct {
	called bool
	err    error
}

func (m *mockCheckRepo) PostResult(_ context.Context, _ int64, _ string, _ int, _ analyzepr.Output) error {
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
	okResult := analyzepr.Output{
		RiskLevel:        entity.RiskLevelLow,
		PriorityScore:    1,
		EstimatedMinutes: 5,
	}

	tests := []struct {
		name        string
		pr          *mockPRRepo
		analyzer    *mockAnalyzerUC
		check       *mockCheckRepo
		wantErr     bool
		wantChecked bool
	}{
		{
			name:        "success",
			pr:          &mockPRRepo{diff: "diff content"},
			analyzer:    &mockAnalyzerUC{result: okResult},
			check:       &mockCheckRepo{},
			wantErr:     false,
			wantChecked: true,
		},
		{
			name:        "pr repo error",
			pr:          &mockPRRepo{err: errors.New("github error")},
			analyzer:    &mockAnalyzerUC{},
			check:       &mockCheckRepo{},
			wantErr:     true,
			wantChecked: false,
		},
		{
			name:        "analyzer error",
			pr:          &mockPRRepo{diff: "diff content"},
			analyzer:    &mockAnalyzerUC{err: errors.New("gemini error")},
			check:       &mockCheckRepo{},
			wantErr:     true,
			wantChecked: false,
		},
		{
			name:        "check repo error",
			pr:          &mockPRRepo{diff: "diff content"},
			analyzer:    &mockAnalyzerUC{result: okResult},
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
