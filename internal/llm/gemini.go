package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kidskoding/music-agent/internal/agent"
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
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &LLMClient{client: client}, nil
}

func (l *LLMClient) SelectNextTrack(
	ctx context.Context, 
	history []agent.Track, 
	available []agent.Track, 
	currentMood string,
) (*agent.Track, string, error) {
	var trackListBuilder strings.Builder
	for _, t := range available {
		trackListBuilder.WriteString(fmt.Sprintf("%s | %s | %s | %s | %.2f\n", t.ID, t.Title, t.Artist, t.Mood, t.Energy))
	}

	lastPlayed := "None"
	if len(history) > 0 {
		last := history[len(history) - 1]
		lastPlayed = fmt.Sprintf("%s by %s", last.Title, last.Artist)
	}

	prompt := fmt.Sprintf(`
		You are an expert DJ AI. 
		Current Vibe: %s.
		Last Played: %s.

		Here is the list of available tracks in my crate:
		%s

		TASK:
		Select the SINGLE best track ID from the list above to play next.
		You must maintain the flow of the current vibe.
		Do not pick the "Last Played" track.

		RETURN JSON ONLY:
		{
		"track_id": "THE_ID_HERE",
		"reason": "Brief explanation why this fits."
		}
	`, currentMood, lastPlayed, trackListBuilder.String())

	resp, err := l.client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return nil, "", fmt.Errorf("gemini error: %w", err)
	}

	var rawText string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.Text != "" {
				rawText = part.Text
				break
			}
		}
	}

	// parse LLM output
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")

	var selection AISelection
	if err := json.Unmarshal([]byte(rawText), &selection); err != nil {
		return nil, "", fmt.Errorf("failed to parse LLM's JSON response: %w (Raw: %s)", err, rawText)
	}

	for _, t := range available {
		if t.ID == selection.TrackID {
			return &t, selection.Reason, nil
		}
	}

	return nil, "", fmt.Errorf("LLM picked ID %s but it wasn't in the list", selection.TrackID)
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

func (l *LLMClient) Close() {}