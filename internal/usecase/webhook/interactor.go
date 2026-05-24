package webhook

import "context"

type Interactor struct {
	pr       PRRepository
	analyzer AnalyzerRepository
	check    CheckRepository
}

func NewInteractor(pr PRRepository, analyzer AnalyzerRepository, check CheckRepository) *Interactor {
	return &Interactor{pr: pr, analyzer: analyzer, check: check}
}

func (i *Interactor) Execute(ctx context.Context, input Input) error {
	diff, err := i.pr.GetDiff(ctx, input.InstallationID, input.RepoFullName, input.PRNumber)
	if err != nil {
		return err
	}

	result, err := i.analyzer.Analyze(ctx, input.Title, diff)
	if err != nil {
		return err
	}

	return i.check.PostResult(ctx, input.InstallationID, input.RepoFullName, input.PRNumber, result)
}
