package mock

import "context"

type LLMRepository struct{}

func NewLLMRepository() *LLMRepository { return &LLMRepository{} }

func (r *LLMRepository) Generate(_ context.Context, _ string) (string, error) {
	return `{
  "risk_level": "medium",
  "risk_reasons": ["ビジネスロジックの変更", "既存関数の修正"],
  "priority_score": 3,
  "summary": "READMEの挨拶文を変更。影響範囲は限定的。",
  "files": [
    {
      "path": "README.md",
      "risk": "low",
      "category": "docs",
      "note": "挨拶文の変更のみ。ロジックへの影響なし。"
    }
  ],
  "review_focus": [
    "変更内容は軽微で、レビュー優先度は低い",
    "他ファイルへの影響がないか確認"
  ],
  "breaking_changes": [],
  "coverage_drop": [],
  "estimated_minutes": 15,
  "custom_output": ""
}`, nil
}
