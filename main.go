package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/kidskoding/music-agent/internal/llm"
	"github.com/kidskoding/music-agent/internal/store"
	"github.com/kidskoding/music-agent/internal/spotify_api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file!")
	}

	var eventStore store.EventStore

	if os.Getenv("DATABRICKS_TOKEN") != "" {
		eventStore, err = store.NewDatabricksStore()
		if err != nil {
			fmt.Printf("failed to connect to Databricks: %v\n", err)
			return
		} else {
			fmt.Println("connected to Databricks Delta Lake")
		}
	} else {
		fmt.Println("no Databricks credentials found")
		return
	}

	ctx := context.Background()
	spotifyClient, err := spotify_api.NewSpotifyClient(ctx)
	if err != nil {
		log.Fatalf("failed to create Spotify client: %v", err)
	}

	llmClient, err := llm.NewLLMClient(ctx)
	if err != nil {
		log.Printf("gemini warning: %v", err)
	}

	playlistID := os.Getenv("TEST_PLAYLIST_ID")
	if playlistID == "" {
		playlistID = "37i9dQZF1DWZeKCadgRdKQ"
	}

	fmt.Println("fetching tracks from Spotify")
	availableTracks, err := spotifyClient.FetchPlaylistTracks(ctx, playlistID)
	if err != nil {
		log.Fatalf("could not fetch tracks: %v", err)
	}
	fmt.Printf("   -> Loaded %d tracks into the Crate.\n", len(availableTracks))

	if llmClient != nil {
		targetMood := "Focus"
		fmt.Printf("asking Gemini to pick a '%s' track...\n", targetMood)
		
		history := []agent.Track{}

		selected, reason, err := llmClient.SelectNextTrack(ctx, history, availableTracks, targetMood)
		if err != nil {
			fmt.Printf("LLM Error: %v\n", err)
		} else {
			fmt.Println("------------------------------------------------")
			fmt.Printf("🎵 SELECTED: %s - %s\n", selected.Title, selected.Artist)
			fmt.Printf("📊 Vibe: %s (Energy: %.2f)\n", selected.Mood, selected.Energy)
			fmt.Printf("🤖 Reason: %s\n", reason)
			fmt.Println("------------------------------------------------")
		}
	}

	agent.StartAgent(eventStore)
}