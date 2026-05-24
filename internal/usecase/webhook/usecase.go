package webhook

import (
	"context"

	"github.com/Integral-07/prism/internal/domain/entity"
)

type Input struct {
	InstallationID int64
	RepoFullName   string
	PRNumber       int
	Title          string
}

type AnalysisResult struct {
	RiskLevel        entity.RiskLevel
	PriorityScore    int
	EstimatedMinutes int
}

type UseCase interface {
	Execute(ctx context.Context, input Input) error
}

type PRRepository interface {
	GetDiff(ctx context.Context, installationID int64, repoFullName string, prNumber int) (string, error)
}

type AnalyzerRepository interface {
	Analyze(ctx context.Context, title, diff string) (AnalysisResult, error)
}

type CheckRepository interface {
	PostResult(ctx context.Context, installationID int64, repoFullName string, prNumber int, result AnalysisResult) error
}
