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

	currentTrack, err := spotifyClient.GetCurrentlyPlaying(ctx)
	if err != nil {
		log.Fatalf("please ensure you are playing music on Spotify! error: %v", err)
	}
	fmt.Printf("listening to: %s - %s (Energy: %.2f, Mood: %s)\n", 
		currentTrack.Title, currentTrack.Artist, currentTrack.Energy, currentTrack.Mood)

	availableTracks, err := spotifyClient.GetUserTopTracks(ctx)
	if err != nil {
		log.Fatalf("could not fetch top tracks: %v", err)
	}
	fmt.Printf("loaded %d tracks into the crate\n", len(availableTracks))

	if llmClient != nil {
		history := []agent.Track{}
		targetMood := currentTrack.Mood

		if targetMood == "" {
			targetMood = "Neutral"
		}

		selected, reason, err := llmClient.SelectNextTrack(ctx, history, availableTracks, targetMood)
		if err != nil {
			fmt.Printf("LLM Error: %v\n", err)
		} else {
			fmt.Println("------------------------------------------------")
			fmt.Printf("selected %s - %s\n", selected.Title, selected.Artist)
			fmt.Printf("vibe: %s with energy: %.2f)\n", selected.Mood, selected.Energy)
			fmt.Printf("reason: %s\n", reason)
			fmt.Println("------------------------------------------------")

			err := spotifyClient.QueueTrack(ctx, selected.ID)
			if err != nil {
				log.Printf("failed to queue track: %v", err)
			}
		}
	}

	agent.StartAgent(eventStore)
}