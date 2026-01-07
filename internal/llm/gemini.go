package llm

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type LLMClient struct {
	client	*genai.Client
}

func NewLLMClient(ctx context.Context) (*LLMClient, error) {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	region := os.Getenv("GOOGLE_REGION")

	if projectID == "" || region == "" {
		return nil, fmt.Errorf("missing GOOGLE_PROJECT_ID or GOOGLE_REGION")
	}

	cfg := &genai.ClientConfig{
		Project:  projectID,
		Location: region,
		Backend:  genai.BackendVertexAI,
	}

	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vertex client: %w", err)
	}

	return &LLMClient{
		client: client,
	}, nil
}

func (l *LLMClient) GetRecommendation(ctx context.Context, currentMood string, currentEnergy float64) (string, error) {
	prompt := fmt.Sprintf(
		"I am a music agent. The user is currently in a '%s' mood with energy level %.1f/1.0. "+
			"Suggest a brief string describing the VIBE of the next track I should play. "+
			"Keep it under 10 words. Do not suggest specific song titles, just the vibe.",
		currentMood, currentEnergy,
	)

	resp, err := l.client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.Text != "" {
				return part.Text, nil
			}
		}
	}

	return "keep the flow steady", nil
}