package middleware

import (
	"net/http"
	"strconv"

	"streamer/internal/playbackauth"
)

func PlaybackAuthentication(
	auth *playbackauth.Service,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			cookie, err := r.Cookie("stream_token")

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := auth.Validate(cookie.Value)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			movieID, err := strconv.Atoi(
				r.PathValue("movieID"),
			)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.MovieID != movieID {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
