package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	analyzepr "github.com/Integral-07/prism/internal/usecase/analyze_pr"
	"github.com/Integral-07/prism/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

type LabelRepository struct {
	clientFor func(installationID int64) (*gh.Client, error)
}

func NewLabelRepository(appID int64, privateKey []byte) *LabelRepository {
	return &LabelRepository{
		clientFor: func(id int64) (*gh.Client, error) {
			return newClient(appID, id, privateKey)
		},
	}
}

var riskLabels = map[entity.RiskLevel]struct {
	name  string
	color string
}{
	entity.RiskLevelHigh:   {"prism: high-risk", "d73a4a"},
	entity.RiskLevelMedium: {"prism: medium-risk", "e4e669"},
	entity.RiskLevelLow:    {"prism: low-risk", "0e8a16"},
}

var timeLabels = []struct {
	name  string
	color string
}{
	{"prism: ~5min", "0075ca"},
	{"prism: ~15min", "0075ca"},
	{"prism: ~30min", "0075ca"},
	{"prism: ~60min", "0075ca"},
	{"prism: 60min+", "0075ca"},
}

func timeLabelName(minutes int) string {
	switch {
	case minutes <= 10:
		return "prism: ~5min"
	case minutes <= 20:
		return "prism: ~15min"
	case minutes <= 45:
		return "prism: ~30min"
	case minutes <= 90:
		return "prism: ~60min"
	default:
		return "prism: 60min+"
	}
}

func (r *LabelRepository) SyncLabels(ctx context.Context, installationID int64, repoFullName string, prNumber int, result analyzepr.Output) error {
	client, err := r.clientFor(installationID)
	if err != nil {
		return err
	}

	owner, repo := splitFullName(repoFullName)

	// 全 prism: ラベルをリポジトリに作成（既存なら無視）
	allLabels := append(timeLabels, struct{ name, color string }{riskLabels[entity.RiskLevelHigh].name, riskLabels[entity.RiskLevelHigh].color})
	for _, l := range allLabels {
		_ = ensureLabel(ctx, client, owner, repo, l.name, l.color)
	}
	for _, info := range riskLabels {
		_ = ensureLabel(ctx, client, owner, repo, info.name, info.color)
	}

	// 既存の prism: ラベルを PR から除去
	current, _, err := client.Issues.ListLabelsByIssue(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return fmt.Errorf("list labels: %w", err)
	}
	for _, l := range current {
		if strings.HasPrefix(l.GetName(), "prism: ") {
			if _, err := client.Issues.RemoveLabelForIssue(ctx, owner, repo, prNumber, l.GetName()); err != nil {
				return fmt.Errorf("remove label %q: %w", l.GetName(), err)
			}
		}
	}

	// 新しいラベルを付与
	newLabels := []string{
		riskLabels[result.RiskLevel].name,
		timeLabelName(result.EstimatedMinutes),
	}
	_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, prNumber, newLabels)
	return err
}

func ensureLabel(ctx context.Context, client *gh.Client, owner, repo, name, color string) error {
	_, resp, err := client.Issues.CreateLabel(ctx, owner, repo, &gh.Label{
		Name:  gh.Ptr(name),
		Color: gh.Ptr(color),
	})
	if err != nil && (resp == nil || resp.StatusCode != http.StatusUnprocessableEntity) {
		return err
	}
	return nil
}
