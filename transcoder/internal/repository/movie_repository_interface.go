package repository

import "transcoder/internal/models"

type MovieRepository interface {
	GetPendingMovie() (*models.Movie, error)
	UpdateStatus(MovieID int64, status models.MovieStatus) error

	// UpdateStatus(
	// 	movie *models.Movie,
	// 	status models.MovieStatus,
	// ) error
}
