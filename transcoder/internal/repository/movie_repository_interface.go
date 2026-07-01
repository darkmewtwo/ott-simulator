package repository

import "transcoder/internal/models"

type MovieRepository interface {
	GetPendingMovie() (*models.Movie, error)
	UpdateStatus(MovieID int64, status models.MovieStatus) error
	UpdateDuration(MovieID int64, durationSeconds int) error
	UpdateHlsPlaylistPath(MovieID int64, hlsPlaylistPath string) error

	// UpdateStatus(
	// 	movie *models.Movie,
	// 	status models.MovieStatus,
	// ) error
}
