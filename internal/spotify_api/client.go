package spotify_api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	api	*spotify.Client
}

func NewSpotifyClient(ctx context.Context) (*Client, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET")
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get spotify token: %w", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return &Client{api: client}, nil
}

func (c *Client) FetchPlaylistTracks(ctx context.Context, playlistID string) ([]agent.Track, error) {
	id := spotify.ID(playlistID)
	playlistItems, err := c.api.GetPlaylistItems(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	var tracks []agent.Track
	var trackIDs []spotify.ID

	// tempMap := make(map[string]*agent.Track)

	for _, item := range playlistItems.Items {
		if item.Track.Track == nil {
			continue
		}
		
		fullTrack := item.Track.Track
		
		t := agent.Track{
			ID:     string(fullTrack.ID),
			Title:  fullTrack.Name,
			Artist: fullTrack.Artists[0].Name,
		}
		
		tracks = append(tracks, t)
		trackIDs = append(trackIDs, fullTrack.ID)
	}

	if len(trackIDs) > 0 {
		features, err := c.api.GetAudioFeatures(ctx, trackIDs...)
		if err != nil {
			log.Printf("warning: could not fetch audio features: %v", err)
		} else {
			for i, feat := range features {
				if feat == nil {
					continue
				}

				tracks[i].Energy = float64(feat.Energy)
				tracks[i].Mood = mapValenceToMood(feat.Valence, feat.Energy)
			}
		}
	}

	return tracks, nil
}

func mapValenceToMood(valence, energy float32) string {
	v := float64(valence)
	e := float64(energy)

	switch {
	case v > 0.6 && e > 0.6:
		return "Happy/Energetic"
	case v > 0.6 && e <= 0.6:
		return "Chill/Happy"
	case v <= 0.4 && e > 0.6:
		return "Aggressive/Angry"
	case v <= 0.4 && e <= 0.4:
		return "Sad/Melancholic"
	default:
		return "Neutral"
	}
}