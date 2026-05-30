package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gh "github.com/google/go-github/v68/github"
)

func testGHClient(t *testing.T, mux *http.ServeMux) *gh.Client {
	t.Helper()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := gh.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u
	return client
}

func fixedClient(client *gh.Client) func(int64) (*gh.Client, error) {
	return func(_ int64) (*gh.Client, error) { return client, nil }
}

func errorClient(err error) func(int64) (*gh.Client, error) {
	return func(_ int64) (*gh.Client, error) { return nil, err }
}

func TestPRRepository_GetDiff(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		clientFor func(int64) (*gh.Client, error)
		handler   http.HandlerFunc
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "diff --git a/foo.go b/foo.go\n+func Foo() {}")
			},
			want: "diff --git a/foo.go b/foo.go\n+func Foo() {}",
		},
		{
			name:    "client error",
			clientFor: errorClient(errors.New("auth error")),
			wantErr: true,
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal error", http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf func(int64) (*gh.Client, error)
			if tt.clientFor != nil {
				cf = tt.clientFor
			} else {
				mux := http.NewServeMux()
				mux.HandleFunc("/repos/owner/repo/pulls/1", tt.handler)
				mux.HandleFunc("/repos/owner/repo/pulls/1.diff", tt.handler)
				cf = fixedClient(testGHClient(t, mux))
			}

			repo := &PRRepository{clientFor: cf}
			diff, err := repo.GetDiff(ctx, 123, "owner/repo", 1)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && diff != tt.want {
				t.Errorf("want %q, got %q", tt.want, diff)
			}
		})
	}
}
