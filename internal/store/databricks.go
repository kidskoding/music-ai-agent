package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/databricks/databricks-sql-go"
	"github.com/joho/godotenv"
	"github.com/kidskoding/music-agent/internal/events"
)

type DatabricksStore struct {
	db *sql.DB
}

func NewDatabricksStore() (*DatabricksStore, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %v\n", err)
	}

	token := os.Getenv("DATABRICKS_TOKEN")
	host := os.Getenv("DATABRICKS_HOST")
	httpPath := os.Getenv("DATABRICKS_HTTP_PATH")

	if token == "" || host == "" || httpPath == "" {
		return nil, fmt.Errorf("missing required DATABRICKS env vars")
	}

	dsn := fmt.Sprintf("token:%s@%s:443%s", token, host, httpPath)
	
	db, err := sql.Open("databricks", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open databricks connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping databricks: %w", err)
	}

	return &DatabricksStore{db: db}, nil
}

func (s *DatabricksStore) SaveTrackEvent(ctx context.Context, event events.TrackEvent) error {
	query := `
		INSERT INTO agent_events 
		(session_id, track_id, track_name, mood, energy, skipped, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		event.SessionID,
		event.TrackID,
		event.TrackName,
		event.Mood,
		event.Energy,
		event.Skipped,
		event.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to insert event to databricks: %w", err)
	}

	log.Printf("☁️ [DATABRICKS] Inserted Event: %s", event.TrackName)
	return nil
}