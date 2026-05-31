package github

import (
	"context"
	"fmt"
	"strings"

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

func (r *CheckRepository) PostResult(ctx context.Context, installationID int64, repoFullName string, prNumber int, result analyzepr.Output, cfg entity.PrismConfig) error {
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
			Title:   gh.Ptr(buildTitle(result, cfg)),
			Summary: gh.Ptr(buildSummary(result, cfg)),
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

func buildTitle(o analyzepr.Output, cfg entity.PrismConfig) string {
	var parts []string
	if cfg.Output.Triage.RiskLevel {
		parts = append(parts, fmt.Sprintf("%s %s", o.RiskLevel.Emoji(), o.RiskLevel))
	}
	if cfg.Output.Triage.PriorityScore {
		parts = append(parts, fmt.Sprintf("優先度 %d/5", o.PriorityScore))
	}
	if cfg.Output.Triage.EstimatedReviewTime {
		parts = append(parts, fmt.Sprintf("推定%dmin", o.EstimatedMinutes))
	}
	if len(parts) == 0 {
		return "PRism Analysis"
	}
	return strings.Join(parts, " · ")
}

func buildSummary(o analyzepr.Output, cfg entity.PrismConfig) string {
	var sb strings.Builder

	if o.Summary != "" {
		sb.WriteString(o.Summary)
		sb.WriteString("\n\n")
	}

	if cfg.Output.Triage.FilePriorityList && len(o.Files) > 0 {
		sb.WriteString("**📁 ファイル優先順位**\n\n")
		sb.WriteString("| ファイル | リスク | カテゴリ | メモ |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, f := range o.Files {
			sb.WriteString(fmt.Sprintf("| `%s` | %s %s | %s | %s |\n",
				f.Path, f.Risk.Emoji(), f.Risk, f.Category, f.Note))
		}
		sb.WriteString("\n")
	}

	if cfg.Output.Support.ReviewFocus && len(o.ReviewFocus) > 0 {
		sb.WriteString("**📋 重点レビュー箇所**\n\n")
		for _, focus := range o.ReviewFocus {
			sb.WriteString(fmt.Sprintf("- %s\n", focus))
		}
		sb.WriteString("\n")
	}

	if o.CustomOutput != "" {
		sb.WriteString("**💬 カスタム分析**\n\n")
		sb.WriteString(o.CustomOutput)
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String())
}
