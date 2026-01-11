// test commit

package spotify_api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kidskoding/music-agent/internal/agent"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Client struct {
	api *spotify.Client
}

const redirectURI = "http://127.0.0.1:8888/callback"

func NewSpotifyClient(ctx context.Context) (*Client, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET")
	}

	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopeUserTopRead,
		),
	)

	ch := make(chan *spotify.Client)

	completeAuth := func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), "random-state-string", r)
		if err != nil {
			http.Error(w, "couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}

		client := spotify.New(auth.Client(r.Context(), token))
		fmt.Fprintf(w, "login completed! you can close this tab")
		ch <- client
	}

	http.HandleFunc("/callback", completeAuth)

	go func() {
		fmt.Println("please log in to Spotify!: ", auth.AuthURL("random-state-string"))
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	client := <-ch
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

func (c *Client) GetUserTopTracks(ctx context.Context) ([]agent.Track, error) {
	limit := 50
	topTracks, err := c.api.CurrentUsersTopTracks(ctx, spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	var tracks []agent.Track
	var trackIDs []spotify.ID

	for _, fullTrack := range topTracks.Tracks {
		t := agent.Track{
			ID:     string(fullTrack.ID),
			Title:  fullTrack.Name,
			Artist: fullTrack.Artists[0].Name,
		}
		tracks = append(tracks, t)
		trackIDs = append(trackIDs, fullTrack.ID)
	}

	features, err := c.api.GetAudioFeatures(ctx, trackIDs...)
	if err == nil {
		for i, feat := range features {
			if feat != nil {
				tracks[i].Energy = float64(feat.Energy)
				tracks[i].Mood = mapValenceToMood(feat.Valence, feat.Energy)
			}
		}
	}

	return tracks, nil
}

func (c *Client) GetCurrentlyPlaying(ctx context.Context) (*agent.Track, error) {
	currentlyPlaying, err := c.api.PlayerCurrentlyPlaying(ctx)
	
	if err != nil {
		return nil, err
	}

	if currentlyPlaying.Item == nil {
		return nil, fmt.Errorf("nothing playing")
	}

	fullTrack := currentlyPlaying.Item
	return &agent.Track{
		ID:     string(fullTrack.ID),
		Title:  fullTrack.Name,
		Artist: fullTrack.Artists[0].Name,
	}, nil
}

func (c *Client) QueueTrack(ctx context.Context, trackID string) error {
    uri := spotify.ID(trackID)
    return c.api.QueueSong(ctx, uri)
}