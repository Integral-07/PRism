package webhook

import (
	"encoding/json"
	"net/http"

	webhookUC "github.com/Integral-07/prism/internal/usecase/webhook"
)

type Handler struct {
	uc webhookUC.UseCase
}

func NewHandler(uc webhookUC.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-GitHub-Event") {
	case "pull_request":
		h.handlePullRequest(w, r)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) handlePullRequest(w http.ResponseWriter, r *http.Request) {
	var event PullRequestEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if event.Action != "opened" && event.Action != "synchronize" {
		w.WriteHeader(http.StatusOK)
		return
	}

	input := webhookUC.Input{
		InstallationID: event.Installation.ID,
		RepoFullName:   event.Repository.FullName,
		PRNumber:       event.PullRequest.Number,
		Title:          event.PullRequest.Title,
	}

	if err := h.uc.Execute(r.Context(), input); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
