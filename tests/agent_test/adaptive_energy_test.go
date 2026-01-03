package agent_test

import (
	"testing"
	"github.com/kidskoding/music-agent/internal/agent"
)

func TestUpdateMoodBasedOnEnergy(t *testing.T) {
	memory := &agent.SessionMemory{
		EnergyHistory: []float64{},
		CurrentMode:   "medium",
	}

	memory.EnergyHistory = []float64{0.1, 0.2, 0.3}
	agent.UpdateMoodBasedOnEnergy(memory)
	if memory.CurrentMode != "chill" {
		t.Errorf("expected 'chill', got %s", memory.CurrentMode)
	}

	memory.EnergyHistory = []float64{0.4, 0.5, 0.6}
	agent.UpdateMoodBasedOnEnergy(memory)
	if memory.CurrentMode != "medium" {
		t.Errorf("expected 'medium', got %s", memory.CurrentMode)
	}

	memory.EnergyHistory = []float64{0.8, 0.9, 0.85}
	agent.UpdateMoodBasedOnEnergy(memory)
	if memory.CurrentMode != "high" {
		t.Errorf("expected 'high', got %s", memory.CurrentMode)
	}

	memory.EnergyHistory = []float64{0.1, 0.2, 0.3, 0.2}
	agent.UpdateMoodBasedOnEnergy(memory)
	if memory.CurrentMode != "chill" {
		t.Errorf("expected 'chill', got %s", memory.CurrentMode)
	}

	memory.EnergyHistory = []float64{0.8, 0.9}
	agent.UpdateMoodBasedOnEnergy(memory)
	if memory.CurrentMode != "high" {
		t.Errorf("expected 'high', got %s", memory.CurrentMode)
	}
}