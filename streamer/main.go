package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const mediaDir = "/media/movies"
const posterDir = "/media/posters"

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

func streamHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.PathValue("filename"))

	fullPath := filepath.Join(mediaDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}

func posterHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.PathValue("filename"))

	fullPath := filepath.Join(posterDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /stream/{filename}", streamHandler)
	mux.HandleFunc("GET /posters/{filename}", posterHandler)

	log.Println("Streamer listening on :8180")

	err := http.ListenAndServe(":8180", mux)
	if err != nil {
		log.Fatal(err)
	}
}
