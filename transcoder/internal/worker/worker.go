package worker

import (
	"log"
	"time"
	"transcoder/internal/models"
	"transcoder/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type Worker struct {
	repository repository.MovieRepository
}

func New(
	repository repository.MovieRepository,
) *Worker {

	return &Worker{
		repository: repository,
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
	log.Printf("movie, ID=%d title=%s status:%s", movie.ID, movie.Title, movie.Status)
	err = w.repository.UpdateStatus(
		movie.ID,
		models.MovieStatusProcessing,
	)
	if err != nil {
		log.Printf("failed to update status: %v", err)
		return
	}

	log.Printf("Movie %d marked PROCESSING", movie.ID)

	time.Sleep(15 * time.Second)

	err = w.repository.UpdateStatus(
		movie.ID,
		models.MovieStatusPending,
	)
	if err != nil {
		log.Printf("failed to update status: %v", err)
		return
	}

	log.Printf("Movie %d marked READY", movie.ID)

}

func (w *Worker) processMovie(movie *models.Movie) error {
	return nil
}

func (w *Worker) Start() {

	for {

		w.processPendingMovie()

		time.Sleep(5 * time.Second)
	}
}
