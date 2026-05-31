package webhook

import (
	"context"

	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
)

type Interactor struct {
	pr       PRRepository
	analyzer AnalyzerUseCase
	check    CheckRepository
	label    LabelRepository
	config   ConfigRepository
}

func NewInteractor(pr PRRepository, analyzer AnalyzerUseCase, check CheckRepository, label LabelRepository, config ConfigRepository) *Interactor {
	return &Interactor{pr: pr, analyzer: analyzer, check: check, label: label, config: config}
}

func (i *Interactor) Execute(ctx context.Context, input Input) error {
	diff, err := i.pr.GetDiff(ctx, input.InstallationID, input.RepoFullName, input.PRNumber)
	if err != nil {
		return err
	}

	cfg, err := i.config.Get(ctx, input.InstallationID, input.RepoFullName)
	if err != nil {
		return err
	}

	result, err := i.analyzer.Execute(ctx, analyzepr.Input{
		Title:        input.Title,
		Diff:         diff,
		CustomPrompt: cfg.Custom,
	})
	if err != nil {
		return err
	}

	if err := i.check.PostResult(ctx, input.InstallationID, input.RepoFullName, input.PRNumber, result, cfg); err != nil {
		return err
	}

	return i.label.SyncLabels(ctx, input.InstallationID, input.RepoFullName, input.PRNumber, result)
}
