package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	webhookHandler "github.com/Integral-07/prism/internal/handler/webhook"
	webhookUC "github.com/Integral-07/prism/internal/usecase/webhook"
)

func main() {
	_ = godotenv.Load()

	// TODO: インフラ層実装後に差し替える
	uc := webhookUC.NewInteractor(nil, nil, nil)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Post("/webhook", webhookHandler.NewHandler(uc).Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
