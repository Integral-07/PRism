package analyzepr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Integral-07/prism/internal/domain/entity"
)

type Interactor struct {
	llm LLMRepository
}

func NewInteractor(llm LLMRepository) *Interactor {
	return &Interactor{llm: llm}
}

func (i *Interactor) Execute(ctx context.Context, input Input) (Output, error) {
	resp, err := i.llm.Generate(ctx, buildPrompt(input))
	if err != nil {
		return Output{}, err
	}
	return parseResponse(resp)
}

func buildPrompt(input Input) string {
	return fmt.Sprintf(`以下のPull Requestを分析し、JSON形式のみで結果を返してください。

PRタイトル: %s

diff:
%s

返却形式（他のテキストは含めないこと）:
{
  "risk_level": "high" | "medium" | "low",
  "priority_score": 1〜5の整数,
  "estimated_minutes": レビューの推定所要時間（分）
}`, input.Title, input.Diff)
}

type llmResponse struct {
	RiskLevel        string `json:"risk_level"`
	PriorityScore    int    `json:"priority_score"`
	EstimatedMinutes int    `json:"estimated_minutes"`
}

func parseResponse(resp string) (Output, error) {
	var r llmResponse
	if err := json.Unmarshal([]byte(resp), &r); err != nil {
		return Output{}, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	riskLevel := entity.RiskLevel(r.RiskLevel)
	if !riskLevel.IsValid() {
		return Output{}, fmt.Errorf("invalid risk_level: %s", r.RiskLevel)
	}

	return Output{
		RiskLevel:        riskLevel,
		PriorityScore:    r.PriorityScore,
		EstimatedMinutes: r.EstimatedMinutes,
	}, nil
}
