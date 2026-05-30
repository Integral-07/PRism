package github

import (
	"context"

	gh "github.com/google/go-github/v68/github"
)

type PRRepository struct {
	clientFor func(installationID int64) (*gh.Client, error)
}

func NewPRRepository(appID int64, privateKey []byte) *PRRepository {
	return &PRRepository{
		clientFor: func(id int64) (*gh.Client, error) {
			return newClient(appID, id, privateKey)
		},
	}
}

func (r *PRRepository) GetDiff(ctx context.Context, installationID int64, repoFullName string, prNumber int) (string, error) {
	client, err := r.clientFor(installationID)
	if err != nil {
		return "", err
	}

	owner, repo := splitFullName(repoFullName)
	diff, _, err := client.PullRequests.GetRaw(ctx, owner, repo, prNumber, gh.RawOptions{Type: gh.Diff})
	return diff, err
}
