package agent_test

import (
	"context"
	"testing"

	"github.com/kidskoding/music-agent/internal/agent"
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

	_, err := agent.RunAgentStep(context.Background(), testStore, memory, tracks, sessionID)
	if err != nil {
		t.Fatalf("agent step failed: %v", err)
	}

	if len(testStore.Events) != 1 {
		t.Errorf("Expected 1 event in store, got %d", len(testStore.Events))
	}

	savedEvent := testStore.Events[0]
	
	if savedEvent.SessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, savedEvent.SessionID)
	}
	if savedEvent.Mood != "chill" {
		t.Errorf("Expected event mood 'chill', got %s", savedEvent.Mood)
	}
	if savedEvent.TrackName == "" {
		t.Error("Saved event has empty TrackName")
	}
}