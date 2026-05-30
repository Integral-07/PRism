package analyzepr

import (
	"context"

	"github.com/Integral-07/prism/internal/domain/entity"
)

type Input struct {
	Title string
	Diff  string
}

type Output struct {
	RiskLevel        entity.RiskLevel
	PriorityScore    int
	EstimatedMinutes int
}

type UseCase interface {
	Execute(ctx context.Context, input Input) (Output, error)
}

type LLMRepository interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
