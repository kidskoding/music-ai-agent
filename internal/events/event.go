package events

import "time"

type TrackEvent struct {
    SessionID   string      `json:"session_id"`
    TrackID     string      `json:"track_id"`
    TrackName   string      `json:"track_name"`
    Mood        string      `json:"mood"`
    Energy      float64     `json:"energy"`
    Skipped     bool        `json:"skipped"`
    Reason      string      `json:"reason"`
    Timestamp   time.Time   `json:"timestamp"`
}