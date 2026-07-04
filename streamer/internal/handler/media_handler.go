package handler

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"streamer/internal/service"
)

type MediaHandler struct {
	service *service.MediaService
}

func NewMediaHandler(service *service.MediaService) *MediaHandler {
	return &MediaHandler{
		service: service,
	}
}

func (h *MediaHandler) StreamMovie(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.PathValue("filename"))

	fullPath, err := h.service.MoviePath(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return

	}
	http.ServeFile(w, r, fullPath)
}

func (h *MediaHandler) StreamPoster(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.PathValue("filename"))

	fullPath, err := h.service.PosterPath(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return

	}

	http.ServeFile(w, r, fullPath)
}

func (h *MediaHandler) StreamHLS(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.PathValue("movieID")

	movieID, err := strconv.ParseInt(
		movieIDStr,
		10,
		64,
	)
	if err != nil || movieID <= 0 {
		http.Error(
			w,
			"invalid movie id",
			http.StatusBadRequest,
		)
		return
	}

	// if movieID <= 0 {
	// 	http.Error(
	// 		w,
	// 		"invalid movie id",
	// 		http.StatusBadRequest,
	// 	)
	// 	return
	// }

	filename := filepath.Base(r.PathValue("filename"))

	fullPath, err := h.service.HLSPath(
		movieID,
		filename,
	)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return

	}
	contentType := mime.TypeByExtension(
		filepath.Ext(filename),
	)

	switch filepath.Ext(filename) {
	case ".m3u8":
		contentType = "application/vnd.apple.mpegurl"

	case ".ts":
		contentType = "video/mp2t"
	}

	if contentType != "" {
		w.Header().Set(
			"Content-Type",
			contentType,
		)
	}

	http.ServeFile(w, r, fullPath)
}
