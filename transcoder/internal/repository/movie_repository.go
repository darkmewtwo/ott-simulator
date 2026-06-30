package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"transcoder/internal/models"
)

type PostgresMovieRepository struct {
	pool *pgxpool.Pool
}

var _ MovieRepository = (*PostgresMovieRepository)(nil)

func NewPostgresMovieRepository(
	pool *pgxpool.Pool,
) *PostgresMovieRepository {

	return &PostgresMovieRepository{
		pool: pool,
	}
}

func (r *PostgresMovieRepository) GetPendingMovie() (
	*models.Movie,
	error,
) {
	log.Println(r.pool)
	config := r.pool.Config()

	log.Printf("Host: %s", config.ConnConfig.Host)
	log.Printf("User: %s", config.ConnConfig.User)
	log.Printf("Database: %s", config.ConnConfig.Database)
	row := r.pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			title,
			filename,
			status,
			duration_seconds,
			hls_playlist_path
		FROM movies
		WHERE status = 'PENDING'
		LIMIT 1
		`,
	)
	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.FileName,
		&movie.Status,
		&movie.DurationSeconds,
		&movie.HLSPlaylistPath,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &movie, nil

}

func (r *PostgresMovieRepository) UpdateStatus(
	movieID int64,
	status models.MovieStatus,
) error {
	_, err := r.pool.Exec(
		context.Background(),
		`
	UPDATE movies
	SET status = $1
	WHERE id = $2
	`,
		status,
		movieID,
	)
	return err
}
