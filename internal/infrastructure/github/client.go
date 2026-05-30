package github

import (
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v68/github"
)

func newClient(appID, installationID int64, privateKey []byte) (*gh.Client, error) {
	tr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, privateKey)
	if err != nil {
		return nil, err
	}
	return gh.NewClient(&http.Client{Transport: tr}), nil
}

func splitFullName(fullName string) (string, string) {
	parts := strings.SplitN(fullName, "/", 2)
	return parts[0], parts[1]
}
