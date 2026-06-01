package mock

import "context"

type LLMRepository struct{}

func NewLLMRepository() *LLMRepository { return &LLMRepository{} }

func (r *LLMRepository) Generate(_ context.Context, _ string) (string, error) {
	return `{
  "risk_level": "low",
  "risk_reasons": ["ドキュメントのみの変更", "ロジックへの影響なし"],
  "priority_score": 1,
  "summary": "READMEの挨拶文を 'Hello' から 'Hello, PRism!' に変更。影響範囲はドキュメントのみ。",
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
  "estimated_minutes": 5,
  "custom_output": "はい、このPRはREADMEのみの変更です。ソースコードや設定ファイルへの影響は一切ありません。"
}`, nil
}
