package github

import (
	"context"
	"fmt"

	"github.com/Integral-07/prism/internal/domain/entity"
	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
	gh "github.com/google/go-github/v68/github"
)

type CheckRepository struct {
	clientFor func(installationID int64) (*gh.Client, error)
}

func NewCheckRepository(appID int64, privateKey []byte) *CheckRepository {
	return &CheckRepository{
		clientFor: func(id int64) (*gh.Client, error) {
			return newClient(appID, id, privateKey)
		},
	}
}

func (r *CheckRepository) PostResult(ctx context.Context, installationID int64, repoFullName string, prNumber int, result analyzepr.Output) error {
	client, err := r.clientFor(installationID)
	if err != nil {
		return err
	}

	owner, repo := splitFullName(repoFullName)

	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("get PR: %w", err)
	}

	opts := gh.CreateCheckRunOptions{
		Name:       "PRism",
		HeadSHA:    pr.GetHead().GetSHA(),
		Status:     gh.Ptr("completed"),
		Conclusion: gh.Ptr(conclusion(result.RiskLevel)),
		Output: &gh.CheckRunOutput{
			Title:   gh.Ptr("PRism Analysis"),
			Summary: gh.Ptr(formatSummary(result)),
		},
	}

	_, _, err = client.Checks.CreateCheckRun(ctx, owner, repo, opts)
	return err
}

func conclusion(level entity.RiskLevel) string {
	if level == entity.RiskLevelHigh {
		return "neutral"
	}
	return "success"
}

func formatSummary(o analyzepr.Output) string {
	return fmt.Sprintf(
		"**Risk Level**: %s %s\n**Priority Score**: %d / 5\n**Estimated Review Time**: %d min",
		o.RiskLevel.Emoji(), o.RiskLevel,
		o.PriorityScore,
		o.EstimatedMinutes,
	)
}
