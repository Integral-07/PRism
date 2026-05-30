package analyzepr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Integral-07/prism/internal/domain/entity"
	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
)

type mockLLM struct {
	resp string
	err  error
}

func (m *mockLLM) Generate(_ context.Context, _ string) (string, error) {
	return m.resp, m.err
}

func TestInteractor_Execute(t *testing.T) {
	ctx := context.Background()
	input := analyzepr.Input{Title: "feat: add payment", Diff: "+ func Pay() {}"}

	tests := []struct {
		name       string
		llmResp    string
		llmErr     error
		wantOutput analyzepr.Output
		wantErr    bool
	}{
		{
			name:    "success low risk",
			llmResp: `{"risk_level":"low","priority_score":2,"estimated_minutes":10}`,
			wantOutput: analyzepr.Output{
				RiskLevel:        entity.RiskLevelLow,
				PriorityScore:    2,
				EstimatedMinutes: 10,
			},
		},
		{
			name:    "success high risk",
			llmResp: `{"risk_level":"high","priority_score":5,"estimated_minutes":60}`,
			wantOutput: analyzepr.Output{
				RiskLevel:        entity.RiskLevelHigh,
				PriorityScore:    5,
				EstimatedMinutes: 60,
			},
		},
		{
			name:    "llm error",
			llmErr:  errors.New("api error"),
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			llmResp: `not-json`,
			wantErr: true,
		},
		{
			name:    "invalid risk_level",
			llmResp: `{"risk_level":"critical","priority_score":5,"estimated_minutes":30}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := analyzepr.NewInteractor(&mockLLM{resp: tt.llmResp, err: tt.llmErr})
			got, err := uc.Execute(ctx, input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.wantOutput {
				t.Errorf("want %+v, got %+v", tt.wantOutput, got)
			}
		})
	}
}
