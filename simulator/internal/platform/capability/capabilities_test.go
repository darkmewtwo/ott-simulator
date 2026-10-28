package capability_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"simulator/internal/platform/capability/authentication"
	"simulator/internal/platform/capability/catalog"
	"simulator/internal/platform/capability/playback"
	"simulator/internal/platform/httpclient"
)

func TestAuthenticationLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		if r.Method != "POST" || payload["username"] != "alice" || payload["password"] != "secret" {
			t.Errorf("unexpected request %s %v", r.Method, payload)
		}
		switch r.URL.Path {
		case "/auth/register":
			if payload["email"] != "alice@example.com" {
				t.Error("missing email")
			}
			fmt.Fprint(w, `{"id":7,"username":"alice","email":"alice@example.com"}`)
		case "/auth/login":
			fmt.Fprint(w, `{"access_token":"signed","token_type":"bearer"}`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := httpclient.New(server.URL)
	auth := authentication.New()
	registered, err := auth.Register(client, authentication.RegisterRequest{Username: "alice", Password: "secret", Email: "alice@example.com"})
	if err != nil || registered.ID != 7 || registered.Username != "alice" || registered.Email != "alice@example.com" {
		t.Fatalf("registration = %+v, %v", registered, err)
	}
	loggedIn, err := auth.Login(client, authentication.LoginRequest{Username: "alice", Password: "secret"})
	if err != nil || loggedIn.AccessToken != "signed" || client.AccessToken != "signed" || client.TokenType != "bearer" {
		t.Fatalf("login = %+v, %v", loggedIn, err)
	}
	if err := auth.Logout(client); err != nil {
		t.Fatal(err)
	}
	if client.AccessToken != "" || client.TokenType != "" {
		t.Fatal("logout did not clear credentials")
	}
}

func TestCatalogResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s", r.Method)
		}
		switch r.URL.Path {
		case "/movies":
			fmt.Fprint(w, `[{"id":7,"title":"Film","poster_url":null}]`)
		case "/movies/7":
			fmt.Fprint(w, `{"id":7,"stream_url":"/stream/film.mp4","hls_url":null}`)
		case "/watch_progress/watch-history":
			fmt.Fprint(w, `[{"movie":{"id":7},"progress":{"last_position_seconds":42,"watch_percentage":35,"started_at":"2024-01-02T00:00:00Z","is_completed":false}}]`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, cat := httpclient.New(server.URL), catalog.New()
	movies, err := cat.ListMovies(client)
	if err != nil || len(movies) != 1 {
		t.Fatalf("movies = %+v, %v", movies, err)
	}
	if movies[0].ID != 7 || movies[0].Title != "Film" || movies[0].PosterURL != nil {
		t.Errorf("movie = %+v", movies[0])
	}
	movie, err := cat.GetMovie(client, 7)
	if err != nil || movie.ID != 7 || movie.StreamURL != "/stream/film.mp4" || movie.HLSURL != nil {
		t.Fatalf("details = %+v, %v", movie, err)
	}
	history, err := cat.GetWatchHistory(client)
	if err != nil || len(history) != 1 {
		t.Fatalf("history = %+v, %v", history, err)
	}
	if history[0].Movie.ID != 7 || history[0].Progress.LastPositionSeconds != 42 || history[0].Progress.StartedAt == nil || history[0].Progress.StartedAt.Year() != 2024 {
		t.Errorf("history entry = %+v", history[0])
	}
}

func TestPlayback(t *testing.T) {
	for _, body := range []string{`null`, `{"movie_id":7,"last_position_seconds":42,"is_completed":true}`} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/events":
					var event playback.WatchEventRequest
					if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
						t.Error(err)
					}
					if r.Method != "POST" || event.MovieID != 7 || event.EventType != playback.EventHeartbeat || event.PositionSeconds != 42 {
						t.Errorf("event = %+v, method = %s", event, r.Method)
					}
					fmt.Fprint(w, `{"id":1,"user_id":3,"movie_id":7,"event_type":"HEARTBEAT","position_seconds":42}`)
				case "/watch_progress/movie/7":
					if r.Method != "GET" {
						t.Errorf("method = %s", r.Method)
					}
					fmt.Fprint(w, body)
				default:
					t.Errorf("unexpected path %s", r.URL.Path)
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			client, player := httpclient.New(server.URL), playback.New()
			event, err := player.CreateEvent(client, playback.WatchEventRequest{MovieID: 7, EventType: playback.EventHeartbeat, PositionSeconds: 42})
			if err != nil || event.ID != 1 || event.MovieID != 7 || event.EventType != playback.EventHeartbeat || event.PositionSeconds != 42 {
				t.Fatalf("event = %+v, %v", event, err)
			}
			progress, err := player.GetMovieProgress(client, 7)
			if err != nil {
				t.Fatal(err)
			}
			if body == "null" {
				if progress != nil {
					t.Errorf("unwatched progress = %+v", progress)
				}
			} else if progress == nil || progress.MovieID != 7 || progress.LastPositionSeconds != 42 || !progress.IsCompleted {
				t.Errorf("progress = %+v", progress)
			}
		})
	}
}

func TestCapabilitiesPropagateHTTPFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "failed", 500) }))
	defer server.Close()
	client := httpclient.New(server.URL)
	client.AccessToken, client.TokenType = "existing", "bearer"
	auth, cat, player := authentication.New(), catalog.New(), playback.New()
	if out, err := auth.Login(client, authentication.LoginRequest{}); out != nil || err == nil {
		t.Errorf("login = %+v, %v", out, err)
	}
	if client.AccessToken != "existing" {
		t.Error("failed login overwrote token")
	}
	if out, err := auth.Register(client, authentication.RegisterRequest{}); out != nil || err == nil {
		t.Errorf("register = %+v, %v", out, err)
	}
	if out, err := cat.ListMovies(client); out != nil || err == nil {
		t.Errorf("list = %+v, %v", out, err)
	}
	if out, err := cat.GetMovie(client, 7); out != nil || err == nil {
		t.Errorf("movie = %+v, %v", out, err)
	}
	if out, err := cat.GetWatchHistory(client); out != nil || err == nil {
		t.Errorf("history = %+v, %v", out, err)
	}
	if out, err := player.CreateEvent(client, playback.WatchEventRequest{}); out != nil || err == nil {
		t.Errorf("event = %+v, %v", out, err)
	}
	if out, err := player.GetMovieProgress(client, 7); out != nil || err == nil {
		t.Errorf("progress = %+v, %v", out, err)
	}
}
