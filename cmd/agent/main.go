package main

import (
	"fmt"

	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/kidskoding/music-agent/internal/store"
)

func main() {
	fmt.Println("starting music agent!")

	eventStore := store.NewLocalStore()
	agent.StartAgent(eventStore)
}