package playback

import (
	"fmt"
	"time"

	"simulator/internal/platform/httpclient"
)

type EventType string

const (
	EventPlay      EventType = "PLAY"
	EventPause     EventType = "PAUSE"
	EventSeek      EventType = "SEEK"
	EventStop      EventType = "STOP"
	EventComplete  EventType = "COMPLETE"
	EventHeartbeat EventType = "HEARTBEAT"
)

type WatchEventRequest struct {
	MovieID         int       `json:"movie_id"`
	EventType       EventType `json:"event_type"`
	PositionSeconds int       `json:"position_seconds"`
}

type WatchEventResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	MovieID         int       `json:"movie_id"`
	EventType       EventType `json:"event_type"`
	PositionSeconds int       `json:"position_seconds"`
}

type WatchProgressResponse struct {
	MovieID             int  `json:"movie_id"`
	LastPositionSeconds int  `json:"last_position_seconds"`
	IsCompleted         bool `json:"is_completed"`
}

type ContinueWatchingResponse struct {
	MovieID             int     `json:"movie_id"`
	Title               string  `json:"title"`
	PosterURL           *string `json:"poster_url"`
	LastPositionSeconds int     `json:"last_position_seconds"`
}

type WatchHistoryMovieResponse struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	PosterURL       *string `json:"poster_url"`
	DurationSeconds int     `json:"duration_seconds"`
}

type WatchHistoryProgressResponse struct {
	LastPositionSeconds int        `json:"last_position_seconds"`
	WatchPercentage     float64    `json:"watch_percentage"`
	StartedAt           *time.Time `json:"started_at"`
	LastWatchedAt       *time.Time `json:"last_watched_at"`
	IsCompleted         bool       `json:"is_completed"`
}

type WatchHistoryResponse struct {
	Movie    WatchHistoryMovieResponse    `json:"movie"`
	Progress WatchHistoryProgressResponse `json:"progress"`
}

func (c *Capability) CreateEvent(
	client *httpclient.Client,
	request WatchEventRequest,
) (*WatchEventResponse, error) {

	var response WatchEventResponse

	err := client.Post("/events", request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Capability) GetMovieProgress(
	client *httpclient.Client,
	movieID int,
) (*WatchProgressResponse, error) {

	var response *WatchProgressResponse

	path := fmt.Sprintf(
		"/watch_progress/movie/%d",
		movieID,
	)

	err := client.Get(path, &response)
	if err != nil {
		return nil, err
	}

	// A nil response means the user has not watched this movie.
	return response, nil
}

func (c *Capability) GetProgress(
	client *httpclient.Client,
) ([]WatchProgressResponse, error) {

	var response []WatchProgressResponse

	err := client.Get(
		"/watch_progress/progress",
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Capability) GetContinueWatching(
	client *httpclient.Client,
) ([]ContinueWatchingResponse, error) {

	var response []ContinueWatchingResponse

	err := client.Get(
		"/watch_progress/continue-watching",
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Capability) GetWatchHistory(
	client *httpclient.Client,
) ([]WatchHistoryResponse, error) {

	var response []WatchHistoryResponse

	err := client.Get(
		"/watch_progress/watch-history",
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
