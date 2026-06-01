package analyzepr

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`以下のPull Requestを分析し、JSON形式のみで結果を返してください。

PRタイトル: %s

diff:
%s

返却形式（他のテキストは含めないこと）:
{
  "risk_level": "high" | "medium" | "low",
  "risk_reasons": ["理由1", "理由2"],
  "priority_score": 1〜5の整数,
  "summary": "このPRの要約（1〜2文）",
  "files": [
    {
      "path": "ファイルパス",
      "risk": "high" | "medium" | "low",
      "category": "logic" | "test" | "docs" | "config" | "style",
      "note": "このファイルで注意すべき点"
    }
  ],
  "review_focus": ["注意すべき点1", "注意すべき点2"],
  "breaking_changes": ["破壊的変更1", "破壊的変更2"],
  "coverage_drop": ["カバレッジ低下が懸念される箇所1"],
  "estimated_minutes": レビューの推定所要時間（分）,
  "custom_output": ""
}`, input.Title, input.Diff))

	if input.CustomPrompt != "" {
		sb.WriteString(fmt.Sprintf("\n\n追加指示:\n%s\ncustom_outputフィールドにその回答を記載してください。", input.CustomPrompt))
	}

	return sb.String()
}

type llmFileRisk struct {
	Path     string `json:"path"`
	Risk     string `json:"risk"`
	Category string `json:"category"`
	Note     string `json:"note"`
}

type llmResponse struct {
	RiskLevel        string        `json:"risk_level"`
	RiskReasons      []string      `json:"risk_reasons"`
	PriorityScore    int           `json:"priority_score"`
	Summary          string        `json:"summary"`
	Files            []llmFileRisk `json:"files"`
	ReviewFocus      []string      `json:"review_focus"`
	BreakingChanges  []string      `json:"breaking_changes"`
	CoverageDrop     []string      `json:"coverage_drop"`
	EstimatedMinutes int           `json:"estimated_minutes"`
	CustomOutput     string        `json:"custom_output"`
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

	files := make([]FileRisk, 0, len(r.Files))
	for _, f := range r.Files {
		risk := entity.RiskLevel(f.Risk)
		if !risk.IsValid() {
			risk = entity.RiskLevelLow
		}
		files = append(files, FileRisk{
			Path:     f.Path,
			Risk:     risk,
			Category: f.Category,
			Note:     f.Note,
		})
	}

	return Output{
		RiskLevel:        riskLevel,
		RiskReasons:      r.RiskReasons,
		PriorityScore:    r.PriorityScore,
		Summary:          r.Summary,
		Files:            files,
		ReviewFocus:      r.ReviewFocus,
		BreakingChanges:  r.BreakingChanges,
		CoverageDrop:     r.CoverageDrop,
		EstimatedMinutes: r.EstimatedMinutes,
		CustomOutput:     r.CustomOutput,
	}, nil
}
