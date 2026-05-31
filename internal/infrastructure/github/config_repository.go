package github

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/Integral-07/prism/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
	"gopkg.in/yaml.v3"
)

type ConfigRepository struct {
	clientFor func(installationID int64) (*gh.Client, error)
}

func NewConfigRepository(appID int64, privateKey []byte) *ConfigRepository {
	return &ConfigRepository{
		clientFor: func(id int64) (*gh.Client, error) {
			return newClient(appID, id, privateKey)
		},
	}
}

func (r *ConfigRepository) Get(ctx context.Context, installationID int64, repoFullName string) (entity.PrismConfig, error) {
	client, err := r.clientFor(installationID)
	if err != nil {
		return entity.DefaultPrismConfig(), err
	}

	owner, repo := splitFullName(repoFullName)

	file, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, ".prism.yml", nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return entity.DefaultPrismConfig(), nil
		}
		return entity.DefaultPrismConfig(), err
	}

	content, err := file.GetContent()
	if err != nil {
		return entity.DefaultPrismConfig(), err
	}
	raw, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return entity.DefaultPrismConfig(), err
	}

	cfg := entity.DefaultPrismConfig()
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return entity.DefaultPrismConfig(), err
	}

	return cfg, nil
}
