package agent

import (
	"fmt"
	"time"
)

func StartAgent() {
	memory := &SessionMemory{
		LastTracks: []string{},
		SkipHistory: make(map[string]bool),
		EnergyHistory: []float64{},
		CurrentMode: "medium",
	}

	for i := 0; i < 3; i++ {
		fmt.Println("agent decides next track...")
		DecideNextTrack(memory)
		time.Sleep(1 * time.Second)
	}
}

func DecideNextTrack(memory *SessionMemory) {
	nextTrack := "track_" + fmt.Sprint(len(memory.LastTracks)+1)
    memory.LastTracks = append(memory.LastTracks, nextTrack)
    fmt.Println("Next track selected:", nextTrack)
}