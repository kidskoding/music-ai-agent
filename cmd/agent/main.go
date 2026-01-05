package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kidskoding/music-agent/internal/agent"
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

	agent.StartAgent(eventStore)
}