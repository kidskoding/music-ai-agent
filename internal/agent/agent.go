package agent

import (
	"fmt"
	"math/rand"
	"time"
)

var SampleTracksExport = []*Track{
	{ID: "1", Title: "Morning Breeze", Mood: "chill", Energy: 0.2, Genre: "ambient"},
	{ID: "2", Title: "Afternoon Drive", Mood: "medium", Energy: 0.5, Genre: "pop"},
	{ID: "3", Title: "Night Run", Mood: "high", Energy: 0.8, Genre: "electronic"},
	{ID: "4", Title: "Evening Calm", Mood: "chill", Energy: 0.3, Genre: "jazz"},
	{ID: "5", Title: "Party Time", Mood: "high", Energy: 0.9, Genre: "dance"},
	{ID: "6", Title: "Lazy Afternoon", Mood: "medium", Energy: 0.4, Genre: "indie"},
	{ID: "7", Title: "Sunset Chill", Mood: "chill", Energy: 0.25, Genre: "ambient"},
	{ID: "8", Title: "Drive Fast", Mood: "high", Energy: 0.85, Genre: "electronic"},
}

func DecideNextTrack(memory *SessionMemory, tracks []*Track) *Track {
	if len(tracks) == 0 {
		fmt.Println("no tracks provided to agent")
		return nil
	}

	var candidates []*Track
	for _, track := range tracks {
		if track.Mood == memory.CurrentMode && 
			!wasRecentlyPlayed(memory, track.ID) &&
			!wasSkipped(memory, track.ID) {
			candidates = append(candidates, track)
		}
	}

	if len(candidates) == 0 {
		fmt.Println("no tracks match the current mood:", memory.CurrentMode)
		return nil
	}

	var nextTrack *Track
	if len(candidates) > 0 {
		rand.Seed(time.Now().UnixNano())
		nextTrack = candidates[0]
	} else {
		for _, track := range tracks {
			if !wasSkipped(memory, track.ID) {
				nextTrack = track
				break
			}
		}

		if nextTrack == nil {
			nextTrack = tracks[0]
		}
	}

	memory.LastTracks = append(memory.LastTracks, nextTrack)

	fmt.Printf("next track selected: %s (mood: %s, energy: %.1f)\n",
        nextTrack.Title, nextTrack.Mood, nextTrack.Energy)

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

func UpdateMoodBasedOnEnergy(memory *SessionMemory) {
	history := memory.EnergyHistory
	n := len(history)
	if n == 0 {
		return
	}

	window := 3
	if n < 3 {
		window = n
	}

	var sum float64 = 0.0
	for i := n - window; i < n; i++ {
		sum += history[i]
	}
	avgEnergy := sum / float64(window)

	switch {
	case avgEnergy <= 0.3:
		memory.CurrentMode = "chill"
	case avgEnergy <= 0.7:
		memory.CurrentMode = "medium"
	default:
		memory.CurrentMode = "high"
	}
}