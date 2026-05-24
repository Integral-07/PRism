package webhook_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Integral-07/prism/internal/handler/webhook"
	webhookUC "github.com/Integral-07/prism/internal/usecase/webhook"
)

type mockUseCase struct {
	err error
}

func (m *mockUseCase) Execute(_ context.Context, _ webhookUC.Input) error {
	return m.err
}

func newRequest(event, body string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	r.Header.Set("X-GitHub-Event", event)
	r.Header.Set("Content-Type", "application/json")
	return r
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name       string
		event      string
		body       string
		ucErr      error
		wantStatus int
	}{
		{
			name:       "unknown event",
			event:      "push",
			body:       `{}`,
			wantStatus: http.StatusOK,
		},
		{
			name:  "pull_request opened",
			event: "pull_request",
			body: `{
				"action": "opened",
				"pull_request": {"number": 1, "title": "feat: add feature", "html_url": "https://github.com/owner/repo/pull/1"},
				"repository": {"full_name": "owner/repo"},
				"installation": {"id": 123}
			}`,
			wantStatus: http.StatusOK,
		},
		{
			name:  "pull_request synchronize",
			event: "pull_request",
			body: `{
				"action": "synchronize",
				"pull_request": {"number": 1, "title": "feat: add feature", "html_url": "https://github.com/owner/repo/pull/1"},
				"repository": {"full_name": "owner/repo"},
				"installation": {"id": 123}
			}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "pull_request ignored action",
			event:      "pull_request",
			body:       `{"action": "closed", "pull_request": {}, "repository": {}, "installation": {}}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "pull_request invalid JSON",
			event:      "pull_request",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "usecase error",
			event: "pull_request",
			body: `{
				"action": "opened",
				"pull_request": {"number": 1, "title": "feat: add feature"},
				"repository": {"full_name": "owner/repo"},
				"installation": {"id": 123}
			}`,
			ucErr:      errors.New("something went wrong"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mockUseCase{err: tt.ucErr}
			h := webhook.NewHandler(uc)

			w := httptest.NewRecorder()
			h.Handle(w, newRequest(tt.event, tt.body))

			if w.Code != tt.wantStatus {
				t.Errorf("want %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
