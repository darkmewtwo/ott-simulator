package main

import (
	"encoding/json"
	"log"
	"net/http"
	"streamer/internal/handler"
	"streamer/internal/middleware"
	"streamer/internal/service"
)

const mediaDir = "/media/movies/"
const posterDir = "/media/posters/"
const hlsDir = "/media/hls/"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(
		map[string]string{
			"service": "streamer",
			"status":  "defnitely healthy",
		},
	)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{
	// 	"service":"streamer",
	// 	"status":"healthy"
	// }`))
}

func main() {
	mux := http.NewServeMux()

	mediaService := service.NewMediaService(mediaDir, posterDir, hlsDir)

	mediaHandler := handler.NewMediaHandler(mediaService)

	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("GET /stream/{filename}", mediaHandler.StreamMovie)

	mux.HandleFunc("GET /posters/{filename}", mediaHandler.StreamPoster)

	mux.HandleFunc("GET /hls/{movieID}/{filename}", mediaHandler.StreamHLS)

	log.Println("Streamer listening on :8180")

	handlerWithMiddleware := middleware.CORS(
		"http://localhost:3000",
	)(mux)

	err := http.ListenAndServe(":8180", handlerWithMiddleware)
	if err != nil {
		log.Fatal(err)
	}
}
