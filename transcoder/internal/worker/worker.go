package worker

import (
	"log"
	"time"
	"transcoder/internal/ffmpeg"
	"transcoder/internal/models"
	"transcoder/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type Worker struct {
	repository repository.MovieRepository
	ffmpeg     *ffmpeg.Service
}

func New(
	repository repository.MovieRepository,
	ffmpeg *ffmpeg.Service,
) *Worker {

	return &Worker{
		repository: repository,
		ffmpeg:     ffmpeg,
	}
}

func (w *Worker) processPendingMovie() {

	log.Println("Checking for pending movies...")
	movie, err := w.repository.GetPendingMovie()

	if err != nil {
		if ce, ok := err.(*pgconn.ConnectError); ok {
			log.Printf("Host: %s", ce.Config.Host)
			log.Printf("Port: %d", ce.Config.Port)
			log.Printf("User: %s", ce.Config.User)
			log.Printf("Database: %s", ce.Config.Database)
		}

		log.Printf("Encountered ERROR: %#v\n", err)
		return
	}

	if movie == nil {
		log.Println("no pending movies")
		return
	}

	err = w.processMovie(movie)
	if err != nil {
		log.Printf(
			"failed to process movie %d: %v",
			movie.ID,
			err,
		)
	}

}

func (w *Worker) processMovie(movie *models.Movie) error {
	log.Printf(
		"Processing movie %d (%s)",
		movie.ID,
		movie.Title,
	)

	err := w.repository.UpdateStatus(
		movie.ID,
		models.MovieStatusProcessing,
	)
	if err != nil {
		return err
	}

	log.Printf(
		"Movie %d marked PROCESSING",
		movie.ID,
	)
	duration, err := w.ffmpeg.GetDuration(movie.FileName)

	if err != nil {
		return err
	}

	log.Printf("Movie duration: %d seconds",
		duration,
	)

	err = w.repository.UpdateDuration(movie.ID, duration)
	if err != nil {
		return err
	}

	log.Printf(
		"Movie %d duration: %d seconds updated in db",
		movie.ID,
		duration,
	)

	// time.Sleep(10 * time.Second)

	playlistPath, err := w.ffmpeg.GenerateHlsPlaylist(
		movie.ID,
		movie.FileName,
	)

	if err != nil {
		return err
	}
	log.Printf("hls Playlist path: %s", playlistPath)

	err = w.repository.UpdateHlsPlaylistPath(
		movie.ID,
		playlistPath,
	)

	if err != nil {
		return err
	}

	err = w.repository.UpdateStatus(
		movie.ID,
		models.MovieStatusReady,
	)

	if err != nil {
		return err
	}

	log.Printf(
		"Movie %d marked READY",
		movie.ID,
	)

	return nil
}

func (w *Worker) Start() {

	for {

		w.processPendingMovie()

		time.Sleep(5 * time.Second)
	}
	// w.processPendingMovie()
}
