package middleware

import (
	"net/http"
	"streamer/internal/auth"
	"strings"
)

func Authentication(
	authService *auth.AuthService,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			// Ignore health and poster requests,
			// Below code is for learning purpose how http handlers are passed through
			switch {
			case r.URL.Path == "/health":
				next.ServeHTTP(w, r)
				return

			case strings.HasPrefix(r.URL.Path, "/posters/"):
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(
					w,
					"invalid authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			token := strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

			err := authService.ValidateToken(token)

			if err != nil {
				http.Error(
					w,
					"invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
