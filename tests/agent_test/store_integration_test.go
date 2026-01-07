package agent_test

import (
	"context"
	"testing"
	"time"

	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/kidskoding/music-agent/internal/events"
	"github.com/kidskoding/music-agent/internal/store"
)

func TestAgentSavesEventToStore(t *testing.T) {
	testStore := store.NewLocalStore()

	memory := &agent.SessionMemory{
		LastTracks:    []*agent.Track{},
		SkipHistory:   make(map[string]bool),
		EnergyHistory: []float64{},
		CurrentMode:   "chill",
	}

	tracks := agent.SampleTracksExport
	sessionID := "test-session-123"
	
	selected := agent.DecideNextTrack(memory, tracks)
	if selected == nil {
		t.Fatal("Agent failed to select a track")
	}

	event := events.TrackEvent{
		SessionID: sessionID,
		TrackID:   selected.ID,
		TrackName: selected.Title,
		Mood:      selected.Mood,
		Energy:    selected.Energy,
		Skipped:   false,
		Reason:    "Test selection",
		Timestamp: time.Now(),
	}

	err := testStore.LogLocalEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("LogEvent failed: %v", err)
	}
	
	if len(testStore.Events) != 1 {
		t.Errorf("Expected 1 event in store, got %d", len(testStore.Events))
	}

	savedEvent := testStore.Events[0]

	if savedEvent.SessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, savedEvent.SessionID)
	}
	if savedEvent.TrackName == "" {
		t.Error("Saved event has empty TrackName")
	}
}