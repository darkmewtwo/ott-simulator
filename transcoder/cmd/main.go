package main

import (
	"fmt"
	"log"
	"os"
	"transcoder/internal/db"
	"transcoder/internal/ffmpeg"
	"transcoder/internal/repository"
	"transcoder/internal/worker"
)

func main() {
	fmt.Println("HELLO")
	log.Println(
		"transcoder started",
	)

	// for {
	// 	log.Println(
	// 		"polling...",
	// 	)

	// 	time.Sleep(
	// 		5 * time.Second,
	// 	)
	// }

	databaseURL := os.Getenv("DATABASE_URL")
	log.Println(databaseURL)
	pool, err := db.NewPool(databaseURL)
	cfg := pool.Config()

	log.Printf("Host=%s", cfg.ConnConfig.Host)
	log.Printf("User=%s", cfg.ConnConfig.User)
	log.Printf("Database=%s", cfg.ConnConfig.Database)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewPostgresMovieRepository(pool)
	ffmpeg := ffmpeg.New("/media/movies/", "/media/hls/", "ffprobe", "ffmpeg")
	worker := worker.New(repo, ffmpeg)

	worker.Start()

	log.Println(databaseURL)
	defer pool.Close()
	defer log.Println("closed connection")
	log.Println("Connected to PostgreSQL")

}
