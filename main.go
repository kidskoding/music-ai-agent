package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/kidskoding/music-agent/internal/events"
	"github.com/kidskoding/music-agent/internal/llm"
	"github.com/kidskoding/music-agent/internal/spotify_api"
	"github.com/kidskoding/music-agent/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("no .env file found! - error: %v\n", err)
	}

	ctx := context.Background()
	sessionID := fmt.Sprintf("session-%d", time.Now().Unix())

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

	llmClient, err := llm.NewLLMClient(ctx)
	if err != nil {
		log.Printf("gemini warning: %v", err)
	}

	spotifyClient, err := spotify_api.NewSpotifyClient(ctx)
	if err != nil {
		log.Fatalf("failed to create Spotify client: %v", err)
	}

	var lastProcessedTrackID string
	for {
		time.Sleep(15 * time.Second)
		currentTrack, err := spotifyClient.GetCurrentlyPlaying(ctx)
		if err != nil {
			log.Fatalf("please ensure you are playing music on Spotify! error: %v", err)
		}
		fmt.Printf("listening to: %s - %s (Energy: %.2f, Mood: %s)\n", 
			currentTrack.Title, currentTrack.Artist, currentTrack.Energy, currentTrack.Mood)

		if currentTrack.ID == lastProcessedTrackID {
			continue
		}

		availableTracks, err := spotifyClient.GetUserTopTracks(ctx)
		if err != nil {
			log.Fatalf("could not fetch top tracks: %v", err)
		}
		fmt.Printf("loaded %d tracks into the crate\n", len(availableTracks))

		if llmClient != nil {
			history := []agent.Track{*currentTrack}
			targetMood := currentTrack.Mood

			if targetMood == "" {
				targetMood = "Neutral"
			}

			selected, reason, err := llmClient.SelectNextTrack(ctx, history, availableTracks, targetMood)
			if err != nil {
				fmt.Printf("LLM Error: %v\n", err)
				continue
			} else {
				fmt.Println("------------------------------------------------")
				fmt.Printf("selected %s - %s\n", selected.Title, selected.Artist)
				fmt.Printf("vibe: %s with energy: %.2f)\n", selected.Mood, selected.Energy)
				fmt.Printf("reason: %s\n", reason)
				fmt.Println("------------------------------------------------")

				err := spotifyClient.QueueTrack(ctx, selected.ID)
				if err != nil {
					log.Printf("failed to queue track: %v", err)
				} else {
					lastProcessedTrackID = currentTrack.ID

					if eventStore != nil {
						event := events.TrackEvent {
							SessionID:  sessionID,
							TrackID:    selected.ID,
							TrackName:  selected.Title,
							Mood:       selected.Mood,
							Energy:     selected.Energy,
							Skipped:    false,
							Reason:     reason,
							Timestamp:  time.Now(),
						}

						go func() {
							err := eventStore.LogEvent(ctx, event)
							if err != nil {
								log.Printf("failed to log Databricks: %v", err)
							} else {
								fmt.Println("event logged to Databricks")
							}
						}()
					}
				}
			}
		}
	}
}