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
}

func NewInteractor(pr PRRepository, analyzer AnalyzerUseCase, check CheckRepository, label LabelRepository) *Interactor {
	return &Interactor{pr: pr, analyzer: analyzer, check: check, label: label}
}

func (i *Interactor) Execute(ctx context.Context, input Input) error {
	diff, err := i.pr.GetDiff(ctx, input.InstallationID, input.RepoFullName, input.PRNumber)
	if err != nil {
		return err
	}

	result, err := i.analyzer.Execute(ctx, analyzepr.Input{
		Title: input.Title,
		Diff:  diff,
	})
	if err != nil {
		return err
	}

	if err := i.check.PostResult(ctx, input.InstallationID, input.RepoFullName, input.PRNumber, result); err != nil {
		return err
	}

	return i.label.SyncLabels(ctx, input.InstallationID, input.RepoFullName, input.PRNumber, result)
}
