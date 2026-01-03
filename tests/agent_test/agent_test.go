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

func TestDecideNextTrackSkipAndRecent(t *testing.T) {
	memory := &agent.SessionMemory{
		LastTracks:    []*agent.Track{},
		SkipHistory:   map[string]bool{"2": true},
		CurrentMode:   "medium",
		EnergyHistory: []float64{},
	}

	SampleTracks := []*agent.Track{
		{ID: "1", Name: "Morning Breeze", Mood: "chill", Energy: 0.2, Genre: "ambient"},
		{ID: "2", Name: "Afternoon Drive", Mood: "medium", Energy: 0.5, Genre: "pop"},
		{ID: "6", Name: "Lazy Afternoon", Mood: "medium", Energy: 0.4, Genre: "indie"},
	}

	next := agent.DecideNextTrack(memory, SampleTracks)
	if next.ID == "2" {
		t.Errorf("Track 2 is skipped but was selected")
	}

	memory.LastTracks = append(memory.LastTracks, next)

	next2 := agent.DecideNextTrack(memory, SampleTracks)
	if next2.ID == next.ID || next2.ID == "2" {
		t.Errorf("Selected track should not be recently played or skipped, got %s", next2.ID)
	}
}

func TestDecideNextTrackRandomization(t *testing.T) {
	memory := &agent.SessionMemory{
		LastTracks:    []*agent.Track{},
		SkipHistory:   make(map[string]bool),
		CurrentMode:   "chill",
		EnergyHistory: []float64{},
	}

	SampleTracks := []*agent.Track{
		{ID: "1", Name: "Morning Breeze", Mood: "chill", Energy: 0.2, Genre: "ambient"},
		{ID: "4", Name: "Evening Calm", Mood: "chill", Energy: 0.3, Genre: "jazz"},
		{ID: "7", Name: "Sunset Chill", Mood: "chill", Energy: 0.25, Genre: "ambient"},
	}

	selected := make(map[string]bool)
	for i := 0; i < 10; i++ {
		track := agent.DecideNextTrack(memory, SampleTracks)
		if track.Mood != "chill" {
			t.Errorf("Expected mood 'chill', got %s", track.Mood)
		}
		selected[track.ID] = true
	}

	if len(selected) < 2 {
		t.Errorf("Randomization might not be working; only %d unique tracks selected", len(selected))
	}
}