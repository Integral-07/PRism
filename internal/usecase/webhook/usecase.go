package webhook

import (
	"context"

	"github.com/Integral-07/prism/internal/domain/entity"
	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
)

type Input struct {
	InstallationID int64
	RepoFullName   string
	PRNumber       int
	Title          string
}

type UseCase interface {
	Execute(ctx context.Context, input Input) error
}

type PRRepository interface {
	GetDiff(ctx context.Context, installationID int64, repoFullName string, prNumber int) (string, error)
}

type AnalyzerUseCase interface {
	Execute(ctx context.Context, input analyzepr.Input) (analyzepr.Output, error)
}

type CheckRepository interface {
	PostResult(ctx context.Context, installationID int64, repoFullName string, prNumber int, result analyzepr.Output, cfg entity.PrismConfig) error
}

type LabelRepository interface {
	SyncLabels(ctx context.Context, installationID int64, repoFullName string, prNumber int, result analyzepr.Output) error
}

type ConfigRepository interface {
	Get(ctx context.Context, installationID int64, repoFullName string) (entity.PrismConfig, error)
}
