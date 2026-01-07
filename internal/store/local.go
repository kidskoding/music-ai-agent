package store

import (
	"context"
	"fmt"
	"sync"
	
	"github.com/kidskoding/music-agent/internal/events"
)

type LocalStore struct {
	mu		sync.Mutex
	Events	[]events.TrackEvent
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		Events: []events.TrackEvent{},
	}
}

func (s *LocalStore) LogLocalEvent(ctx context.Context, event events.TrackEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Events = append(s.Events, event)

	fmt.Printf("storing saved event: %s (Mood: %s) [Total: %d]\n",
		event.TrackName, event.Mood, len(s.Events))

	return nil
}

func (s *LocalStore) Close() error {
	return nil
}