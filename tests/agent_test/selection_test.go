package agent_test

import (
	"testing"
	"github.com/kidskoding/music-agent/internal/agent"
)

func newMemory(mode string) *agent.SessionMemory {
	return &agent.SessionMemory{
		LastTracks:    []*agent.Track{},
		SkipHistory:   make(map[string]bool),
		CurrentMode:   mode,
		EnergyHistory: []float64{},
	}
}

func TestDecideNextTrack(t *testing.T) {
	memory := newMemory("chill")

	var SampleTracks = []*agent.Track{
		{ID: "1", Title: "Morning Breeze", Artist: "Test Artist A", Mood: "chill", Energy: 0.2},
		{ID: "2", Title: "Afternoon Drive", Artist: "Test Artist B", Mood: "medium", Energy: 0.5},
	}

	selection := agent.DecideNextTrack(memory, SampleTracks)

	if selection == nil {
		t.Fatalf("Expected a track, got nil")
	}

	if selection.Mood != "chill" {
		t.Errorf("Expected mood 'chill', got '%s'", selection.Mood)
	}
}

func TestDecideNextTrackSkipAndRecent(t *testing.T) {
	memory := newMemory("medium")
	memory.SkipHistory["2"] = true

	SampleTracks := []*agent.Track{
		{ID: "2", Title: "Afternoon Drive", Mood: "medium", Energy: 0.5}, // Skipped
		{ID: "6", Title: "Lazy Afternoon", Mood: "medium", Energy: 0.4},  // Valid
	}

	next := agent.DecideNextTrack(memory, SampleTracks)
	if next == nil {
		t.Fatal("Got nil track, expected selection")
	}
	if next.ID == "2" {
		t.Errorf("Track 2 is skipped but was selected")
	}
}

func TestDecideNextTrackRandomization(t *testing.T) {
	memory := newMemory("chill")
	SampleTracks := []*agent.Track{
		{ID: "1", Title: "A", Mood: "chill"},
		{ID: "2", Title: "B", Mood: "chill"},
		{ID: "3", Title: "C", Mood: "chill"},
	}

	selected := make(map[string]bool)

	// Run 20 times to ensure we pick at least 2 different songs
	for i := 0; i < 20; i++ {
		track := agent.DecideNextTrack(memory, SampleTracks)
		if track != nil {
			selected[track.ID] = true
		}
	}

	if len(selected) < 2 {
		t.Errorf("Randomization failure: only picked %d unique tracks out of 3 options", len(selected))
	}
}

func TestAgentWithNoTracks(t *testing.T) {
	memory := newMemory("happy")
	emptyCatalog := []*agent.Track{}

	selection := agent.DecideNextTrack(memory, emptyCatalog)

	if selection != nil {
		t.Errorf("Expected nil for empty catalog, got %v", selection)
	}
}

func TestFallbackWhenNoMoodMatch(t *testing.T) {
	memory := newMemory("super_happy")
	
	SampleTracks := []*agent.Track{
		{ID: "1", Title: "Sad Song", Artist: "Emo Band", Mood: "sad", Energy: 0.1},
	}

	selection := agent.DecideNextTrack(memory, SampleTracks)
	
	// Based on your agent.go logic, it SHOULD fall back to the sad song
	if selection == nil {
		t.Error("Agent failed to fall back to available track")
	}
}