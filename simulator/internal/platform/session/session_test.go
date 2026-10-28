package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"simulator/internal/cognition"
	"simulator/internal/platform/httpclient"
	"simulator/internal/user"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name      string
		register  bool
		failPath  string
		empty     bool
		wantPaths []string
	}{
		{"login and leave", false, "", false, []string{"/auth/login", "/movies"}},
		{"register first", true, "", false, []string{"/auth/register", "/auth/login", "/movies"}},
		{"registration error", true, "/auth/register", false, []string{"/auth/register"}},
		{"login error", false, "/auth/login", false, []string{"/auth/login"}},
		{"catalog error", false, "/movies", false, []string{"/auth/login", "/movies"}},
		{"empty catalog selection", false, "", true, []string{"/auth/login", "/movies"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var paths []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
				if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/register" {
					var payload map[string]string
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Error(err)
					}
					if r.Method != "POST" || payload["username"] != "alice" || payload["password"] != "secret" {
						t.Errorf("unexpected auth request: %s %v", r.Method, payload)
					}
					if r.URL.Path == "/auth/register" && payload["email"] != "alice@example.com" {
						t.Error("registration email missing")
					}
				} else if r.Method != "GET" || r.Header.Get("Authorization") != "bearer signed" {
					t.Error("catalog request missing authentication")
				}
				if r.URL.Path == tc.failPath {
					http.Error(w, "failure", 500)
					return
				}
				switch r.URL.Path {
				case "/auth/login":
					fmt.Fprint(w, `{"access_token":"signed","token_type":"bearer"}`)
				case "/auth/register":
					fmt.Fprint(w, `{"id":1}`)
				case "/movies":
					if tc.empty {
						fmt.Fprint(w, `[]`)
					} else {
						fmt.Fprint(w, `[{"id":7}]`)
					}
				default:
					t.Errorf("unexpected path %s", r.URL.Path)
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			u := &user.User{Identity: user.Identity{Username: "alice", Password: "secret", Email: "alice@example.com"}}
			if tc.empty {
				u.Personality = user.Personality{Impulsiveness: 1, AttentionSpan: 1}
				u.MentalState = user.MentalState{Energy: 1, Happiness: 1, Satisfaction: 1}
			}
			s := NewSession(u, tc.register)
			s.HTTPClient = httpclient.New(server.URL)
			err := s.Run()
			if (err != nil) != (tc.failPath != "") {
				t.Fatalf("Run error = %v", err)
			}
			if !reflect.DeepEqual(paths, tc.wantPaths) {
				t.Errorf("paths = %v, want %v", paths, tc.wantPaths)
			}
			if s.HTTPClient.AccessToken != "" || s.HTTPClient.TokenType != "" {
				t.Error("authentication was not cleared")
			}
		})
	}
}

func TestCatalogActionMapper(t *testing.T) {
	mapper := NewCatalogActionMapper(nil)
	for _, ids := range [][]int{nil, {1, 2}} {
		if got := mapper.Execute(cognition.ActionSelectMovie, ids...); got != nil {
			t.Errorf("invalid selection = %v", got)
		}
	}
	if mapper.Execute("unknown") != nil {
		t.Error("unknown action should do nothing")
	}
	mapper.Register(cognition.ActionSelectMovie, func(ids ...int) []int { return ids })
	if got := mapper.Execute(cognition.ActionSelectMovie, 7); !reflect.DeepEqual(got, []int{7}) {
		t.Errorf("replacement handler = %v", got)
	}
}

func TestCatalogActions(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		switch r.URL.Path {
		case "/movies":
			fmt.Fprint(w, `[{"id":7},{"id":9}]`)
		case "/movies/9":
			fmt.Fprint(w, `{"id":9}`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	s := NewSession(&user.User{MentalState: user.MentalState{Energy: 1}}, false)
	s.HTTPClient = httpclient.New(server.URL)
	mapper := NewCatalogActionMapper(s)
	if got := mapper.Execute(cognition.ActionContinueBrowsing); !reflect.DeepEqual(got, []int{7, 9}) {
		t.Fatalf("movie IDs = %v", got)
	}
	if s.User.MentalState.Energy != .96 {
		t.Errorf("browse energy = %v", s.User.MentalState.Energy)
	}
	before := s.User.MentalState
	mapper.Execute(cognition.ActionSelectMovie, 9)
	if s.User.MentalState.Energy >= before.Energy || s.User.MentalState.Satisfaction <= before.Satisfaction {
		t.Errorf("selection cost not applied: %+v", s.User.MentalState)
	}
	if !reflect.DeepEqual(paths, []string{"/movies", "/movies/9"}) {
		t.Errorf("paths = %v", paths)
	}
}
