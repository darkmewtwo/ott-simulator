package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuthorizationAndDecode(t *testing.T) {
	for _, tc := range []struct{ name, token, kind, want string }{
		{"anonymous", "", "", ""}, {"default bearer", "secret", "", "Bearer secret"}, {"custom type", "secret", "Token", "Token secret"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" || r.URL.Path != "/movies" || r.Header.Get("Authorization") != tc.want {
					t.Errorf("unexpected request: %s %s, auth %q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
				}
				fmt.Fprint(w, `{"id":7}`)
			}))
			defer server.Close()
			client := New(server.URL)
			client.AccessToken, client.TokenType = tc.token, tc.kind
			var out struct{ ID int }
			if err := client.Get("/movies", &out); err != nil || out.ID != 7 {
				t.Fatalf("Get = %+v, %v", out, err)
			}
		})
	}
}

func TestPostEncodesRequestAndDecodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/login" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload["username"] != "alice" {
			t.Errorf("payload = %v, error %v", payload, err)
		}
		fmt.Fprint(w, `{"token":"signed"}`)
	}))
	defer server.Close()
	var out struct{ Token string }
	if err := New(server.URL).Post("/login", map[string]string{"username": "alice"}, &out); err != nil || out.Token != "signed" {
		t.Fatalf("Post = %+v, %v", out, err)
	}
}

func TestResponseFailuresAndIgnoredBody(t *testing.T) {
	for _, method := range []string{"GET", "POST"} {
		for _, tc := range []struct {
			name              string
			status            int
			body              string
			output, wantError bool
		}{
			{"http error", 401, `{}`, true, true}, {"server error", 500, `{}`, true, true},
			{"invalid json", 200, `broken`, true, true}, {"empty json", 200, ``, true, true},
			{"ignored response", 204, ``, false, false},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(tc.status); fmt.Fprint(w, tc.body) }))
				defer server.Close()
				var out any
				if tc.output {
					out = &map[string]any{}
				}
				client := New(server.URL)
				var err error
				if method == "GET" {
					err = client.Get("/", out)
				} else {
					err = client.Post("/", nil, out)
				}
				if (err != nil) != tc.wantError {
					t.Fatalf("error = %v, want error %v", err, tc.wantError)
				}
			})
		}
	}
}

func TestRequestFailures(t *testing.T) {
	client := New(":invalid")
	if err := client.Get("/", nil); err == nil {
		t.Error("expected invalid GET URL error")
	}
	if err := client.Post("/", nil, nil); err == nil {
		t.Error("expected invalid POST URL error")
	}
	if err := client.Post("/", make(chan int), nil); err == nil {
		t.Error("expected JSON encoding error")
	}
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	server.Close()
	client = New(server.URL)
	if err := client.Get("/", nil); err == nil {
		t.Error("expected GET transport error")
	}
	if err := client.Post("/", nil, nil); err == nil {
		t.Error("expected POST transport error")
	}
}
