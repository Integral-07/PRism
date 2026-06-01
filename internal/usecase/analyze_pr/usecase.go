package analyzepr

import (
	"context"

	"github.com/Integral-07/prism/internal/domain/entity"
)

type Input struct {
	Title        string
	Diff         string
	CustomPrompt string
}

type FileRisk struct {
	Path     string
	Risk     entity.RiskLevel
	Category string
	Note     string
}

type Output struct {
	RiskLevel        entity.RiskLevel
	RiskReasons      []string
	PriorityScore    int
	Summary          string
	Files            []FileRisk
	ReviewFocus      []string
	BreakingChanges  []string
	CoverageDrop     []string
	EstimatedMinutes int
	CustomOutput     string
}

type UseCase interface {
	Execute(ctx context.Context, input Input) (Output, error)
}

type LLMRepository interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
