package store

import (
	"context"
	
	"github.com/kidskoding/music-agent/internal/events"
)

type EventStore interface {
	SaveTrackEvent(ctx context.Context, event events.TrackEvent) error
}