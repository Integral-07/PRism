package mock

import "context"

type LLMRepository struct{}

func NewLLMRepository() *LLMRepository { return &LLMRepository{} }

func (r *LLMRepository) Generate(_ context.Context, _ string) (string, error) {
	return `{"risk_level":"medium","priority_score":3,"estimated_minutes":15}`, nil
}
