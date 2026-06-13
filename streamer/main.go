package main

import (
	"encoding/json"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(
		map[string]string{
			"service": "streamer",
			"status":  "defnitely healthy",
		},
	)
}

func main() {
	http.HandleFunc("/health", health)

	http.ListenAndServe(":8180", nil)
}
