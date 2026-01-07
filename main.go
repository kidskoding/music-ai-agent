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
		fmt.Println("connected to Gemini!")

		fmt.Println("asking Gemini for a vibe...")
		suggestion, err := llmClient.GetRecommendation(ctx, "focus", 0.8)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Printf("suggestion: %s\n", suggestion)
		}
	}

	agent.StartAgent(eventStore)
}