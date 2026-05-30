package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type contentGenerator interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

type LLMRepository struct {
	newGenerator func(ctx context.Context) (contentGenerator, func(), error)
}

func NewLLMRepository(apiKey string) *LLMRepository {
	return &LLMRepository{
		newGenerator: func(ctx context.Context) (contentGenerator, func(), error) {
			client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
			if err != nil {
				return nil, nil, fmt.Errorf("gemini client: %w", err)
			}
			model := client.GenerativeModel("gemini-2.5-flash")
			model.GenerationConfig.ResponseMIMEType = "application/json"
			return model, func() { client.Close() }, nil
		},
	}
}

func (r *LLMRepository) Generate(ctx context.Context, prompt string) (string, error) {
	gen, cleanup, err := r.newGenerator(ctx)
	if err != nil {
		return "", err
	}
	defer cleanup()

	resp, err := gen.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generate: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response part type")
	}

	return string(text), nil
}
