package gemini

import (
	"context"
	"errors"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

type mockGenerator struct {
	resp *genai.GenerateContentResponse
	err  error
}

func (m *mockGenerator) GenerateContent(_ context.Context, _ ...genai.Part) (*genai.GenerateContentResponse, error) {
	return m.resp, m.err
}

func TestLLMRepository_Generate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		newGenerator func(context.Context) (contentGenerator, func(), error)
		want         string
		wantErr      bool
	}{
		{
			name: "success",
			newGenerator: fixedGenerator(&mockGenerator{
				resp: &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{Content: &genai.Content{Parts: []genai.Part{genai.Text(`{"risk_level":"low"}`)}}},
					},
				},
			}),
			want: `{"risk_level":"low"}`,
		},
		{
			name:         "generate error",
			newGenerator: fixedGenerator(&mockGenerator{err: errors.New("api error")}),
			wantErr:      true,
		},
		{
			name:         "empty candidates",
			newGenerator: fixedGenerator(&mockGenerator{resp: &genai.GenerateContentResponse{}}),
			wantErr:      true,
		},
		{
			name: "empty parts",
			newGenerator: fixedGenerator(&mockGenerator{
				resp: &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{Content: &genai.Content{Parts: []genai.Part{}}},
					},
				},
			}),
			wantErr: true,
		},
		{
			name: "non-text part",
			newGenerator: fixedGenerator(&mockGenerator{
				resp: &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{Content: &genai.Content{Parts: []genai.Part{genai.Blob{MIMEType: "image/png"}}}},
					},
				},
			}),
			wantErr: true,
		},
		{
			name: "client init error",
			newGenerator: func(_ context.Context) (contentGenerator, func(), error) {
				return nil, nil, errors.New("client init failed")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &LLMRepository{newGenerator: tt.newGenerator}
			got, err := repo.Generate(ctx, "test prompt")
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}

func fixedGenerator(gen contentGenerator) func(context.Context) (contentGenerator, func(), error) {
	return func(_ context.Context) (contentGenerator, func(), error) {
		return gen, func() {}, nil
	}
}
