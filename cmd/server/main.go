package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	webhookHandler "github.com/Integral-07/prism/internal/handler/webhook"
	"github.com/Integral-07/prism/internal/infrastructure/config"
	infraGitHub "github.com/Integral-07/prism/internal/infrastructure/github"
	"github.com/Integral-07/prism/internal/infrastructure/gemini"
	analyzeprUC "github.com/Integral-07/prism/internal/usecase/analyze_pr"
	webhookUC "github.com/Integral-07/prism/internal/usecase/webhook"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	privateKey := []byte(cfg.GitHubPrivateKey)

	prRepo := infraGitHub.NewPRRepository(cfg.GitHubAppID, privateKey)
	checkRepo := infraGitHub.NewCheckRepository(cfg.GitHubAppID, privateKey)
	llmRepo := gemini.NewLLMRepository(cfg.GeminiAPIKey)

	analyzerUC := analyzeprUC.NewInteractor(llmRepo)
	uc := webhookUC.NewInteractor(prRepo, analyzerUC, checkRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Post("/webhook", webhookHandler.NewHandler(uc).Handle)

	log.Printf("server starting on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
