package agent

import (
	"fmt"
	"time"
)

func StartAgent() {
	memory := &SessionMemory{
		LastTracks: []*Track{},
		SkipHistory: make(map[string]bool),
		EnergyHistory: []float64{},
		CurrentMode: "medium",
	}

	var SampleTracks = []*Track{
		{ID: "1", Name: "Morning Breeze", Mood: "chill", Energy: 0.2, Genre: "ambient"},
		{ID: "2", Name: "Afternoon Drive", Mood: "medium", Energy: 0.5, Genre: "pop"},
		{ID: "3", Name: "Night Run", Mood: "high", Energy: 0.8, Genre: "electronic"},
		{ID: "4", Name: "Evening Calm", Mood: "chill", Energy: 0.3, Genre: "jazz"},
		{ID: "5", Name: "Party Time", Mood: "high", Energy: 0.9, Genre: "dance"},
		{ID: "6", Name: "Lazy Afternoon", Mood: "medium", Energy: 0.4, Genre: "indie"},
		{ID: "7", Name: "Sunset Chill", Mood: "chill", Energy: 0.25, Genre: "ambient"},
		{ID: "8", Name: "Drive Fast", Mood: "high", Energy: 0.85, Genre: "electronic"},
	}

	for range 3 {
		fmt.Println("agent decides next track...")
		next := DecideNextTrack(memory, SampleTracks)
		fmt.Printf("played: %s\n\n", next.Name)
		time.Sleep(1 * time.Second)
	}
}

func DecideNextTrack(memory *SessionMemory, tracks []*Track) *Track {
	var candidates []*Track
	for _, track := range tracks {
		if track.Mood == memory.CurrentMode && 
			!wasRecentlyPlayed(memory, track.ID) &&
			!wasSkipped(memory, track.ID) {
			candidates = append(candidates, track)
		}
	}

	var nextTrack *Track
	if len(candidates) > 0 {
		nextTrack = candidates[0]
	} else {
		nextTrack = tracks[0]
	}

	memory.LastTracks = append(memory.LastTracks, nextTrack)

	fmt.Printf("next track selected: %s (mood: %s, energy: %.1f)\n",
        nextTrack.Name, nextTrack.Mood, nextTrack.Energy)

	return nextTrack
}

func wasRecentlyPlayed(memory *SessionMemory, trackID string) bool {
	for _, t := range memory.LastTracks {
		if t.ID == trackID {
			return true
		}
	}

	return false
}

func wasSkipped(memory *SessionMemory, trackID string) bool {
	return memory.SkipHistory[trackID]
}