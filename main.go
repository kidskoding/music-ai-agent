package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/kidskoding/music-agent/internal/llm"
	"github.com/kidskoding/music-agent/internal/store"
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
	llmClient, err := llm.NewLLMClient(ctx)
	if err != nil {
		fmt.Printf("could not connect to AI: %v\n", err)
	} else {
		availableTracks := []agent.Track{
			{ID: "t1", Title: "Weightless", Artist: "Marconi Union", Mood: "Focus", Energy: 0.2},
			{ID: "t2", Title: "Strobe", Artist: "Deadmau5", Mood: "Focus", Energy: 0.8},
			{ID: "t3", Title: "Enter Sandman", Artist: "Metallica", Mood: "Aggressive", Energy: 0.9},
		}

		history := []agent.Track{}

		fmt.Println("asking Gemini for a vibe...")
		selectedTrack, reason, err := llmClient.SelectNextTrack(ctx, history, availableTracks, "Focus")
		if err != nil {
			fmt.Println("LLM error:", err)
		} else {
			fmt.Printf("suggestion: %s by %s\n", selectedTrack.Title, selectedTrack.Artist)
			fmt.Printf("reason: %s\n", reason)
		}
	}

	agent.StartAgent(eventStore)
}