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
        {ID: "3", Title: "Night Run", Artist: "Test Artist C", Mood: "high", Energy: 0.8},
        {ID: "4", Title: "Evening Calm", Artist: "Test Artist A", Mood: "chill", Energy: 0.3},
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
        {ID: "1", Title: "Morning Breeze", Artist: "Test Artist A", Mood: "chill", Energy: 0.2},
        {ID: "2", Title: "Afternoon Drive", Artist: "Test Artist B", Mood: "medium", Energy: 0.5},
        {ID: "6", Title: "Lazy Afternoon", Artist: "Test Artist D", Mood: "medium", Energy: 0.4},
    }

    next := agent.DecideNextTrack(memory, SampleTracks)
    if next == nil {
        t.Fatal("Got nil track, expected selection")
    }
    if next.ID == "2" {
        t.Errorf("Track 2 is skipped but was selected")
    }

    memory.LastTracks = append(memory.LastTracks, next)

    next2 := agent.DecideNextTrack(memory, SampleTracks)
    if next2 != nil && next2.ID == next.ID {
        t.Logf("Agent repeated track %s because no others were available (Acceptable behavior)", next2.ID)
    }
}

func TestDecideNextTrackRandomization(t *testing.T) {
    memory := newMemory("chill")

    SampleTracks := []*agent.Track{
        {ID: "1", Title: "Song A", Artist: "Artist X", Mood: "chill", Energy: 0.2},
        {ID: "2", Title: "Song B", Artist: "Artist Y", Mood: "chill", Energy: 0.3},
        {ID: "3", Title: "Song C", Artist: "Artist Z", Mood: "chill", Energy: 0.25},
    }

    selected := make(map[string]bool)

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
    
    if selection == nil {
        t.Log("Agent returned nil when no mood matched (Valid behavior)")
    } else {
        t.Logf("Agent fell back to %s (Valid behavior)", selection.Title)
    }
}

func TestAllTracksSkipped(t *testing.T) {
    memory := newMemory("happy")
    memory.SkipHistory["1"] = true
    memory.SkipHistory["2"] = true

    tracks := []*agent.Track{
        {ID: "1", Title: "Song A", Artist: "Artist 1", Mood: "happy"},
        {ID: "2", Title: "Song B", Artist: "Artist 2", Mood: "happy"},
    }

    selection := agent.DecideNextTrack(memory, tracks)

    if selection != nil {
        t.Errorf("expected nil when all tracks are skipped, got %v", selection.Title)
    } else {
        t.Log("correctly returned nil when user skipped entire catalog")
    }
}

func TestUnknownMoodHandling(t *testing.T) {
    memory := newMemory("SpicyGarlicButter")

    tracks := []*agent.Track{
        {ID: "1", Title: "Normal Song", Artist: "Band A", Mood: "chill"},
        {ID: "2", Title: "Another Song", Artist: "Band B", Mood: "happy"},
    }

    selection := agent.DecideNextTrack(memory, tracks)

    if selection != nil {
        t.Errorf("expected nil for unknown mood, but agent picked %s", selection.Title)
    } else {
        t.Log("correctly handled unknown mood without crashing")
    }
}