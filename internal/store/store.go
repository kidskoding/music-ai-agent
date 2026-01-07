package store

import (
	"context"
	
	"github.com/kidskoding/music-agent/internal/events"
)

type EventStore interface {
	LogEvent(ctx context.Context, event events.TrackEvent) error
	Close() error
}