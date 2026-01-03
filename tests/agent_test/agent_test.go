package agent_test

import (
    "testing"
    "github.com/kidskoding/music-agent/internal/agent"
)

func TestDecideNextTrack(t *testing.T) {
    memory := &agent.SessionMemory{
        LastTracks:    []*agent.Track{},
        SkipHistory:   make(map[string]bool),
        CurrentMode:   "chill",
        EnergyHistory: []float64{},
    }

	var SampleTracks = []*agent.Track{
		{ID: "1", Name: "Morning Breeze", Mood: "chill", Energy: 0.2, Genre: "ambient"},
		{ID: "2", Name: "Afternoon Drive", Mood: "medium", Energy: 0.5, Genre: "pop"},
		{ID: "3", Name: "Night Run", Mood: "high", Energy: 0.8, Genre: "electronic"},
		{ID: "4", Name: "Evening Calm", Mood: "chill", Energy: 0.3, Genre: "jazz"},
		{ID: "5", Name: "Party Time", Mood: "high", Energy: 0.9, Genre: "dance"},
		{ID: "6", Name: "Lazy Afternoon", Mood: "medium", Energy: 0.4, Genre: "indie"},
		{ID: "7", Name: "Sunset Chill", Mood: "chill", Energy: 0.25, Genre: "ambient"},
		{ID: "8", Name: "Drive Fast", Mood: "high", Energy: 0.85, Genre: "electronic"},
	}

    agent.DecideNextTrack(memory, SampleTracks)

    if len(memory.LastTracks) != 1 {
        t.Errorf("Expected 1 track, got %d", len(memory.LastTracks))
    }

    trackID := memory.LastTracks[0].ID
    found := false
    for _, track := range SampleTracks {
        if track.ID == trackID && track.Mood == "chill" {
            found = true
        }
    }
    if !found {
        t.Errorf("Track selected does not match CurrentMode 'chill'")
    }
}
