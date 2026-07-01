package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// func NewPool(
// 	databaseURL string,
// ) (*pgxpool.Pool, error) {

// 	log.Println(databaseURL)
// 	return pgxpool.New(
// 		context.Background(),
// 		databaseURL,
// 	)
// }

func NewPool(databaseURL string) (*pgxpool.Pool, error) {
	log.Printf("NewPool received: %q", databaseURL)

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	cfg := pool.Config()
	log.Printf("Pool Host=%q User=%q Database=%q",
		cfg.ConnConfig.Host,
		cfg.ConnConfig.User,
		cfg.ConnConfig.Database,
	)

	return pool, nil
}
