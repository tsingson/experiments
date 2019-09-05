package main

import (
	"net/http"
	"os"

	"github.com/jackc/pgx"

	"github.com/rs/zerolog"
)

var (
	log  zerolog.Logger
	pool *pgx.ConnPool
	err  error
)

func init() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""

	log = zerolog.New(os.Stderr)
}

func main() {
	pool, err = InitPgx(log)
	if err != nil {
		log.Panic().Err(err).Msg("Unable to create connection pool")
		os.Exit(1)
	}
	http.HandleFunc("/", urlHandler)

	log.Info().Msg("Starting URL shortener on localhost:8080")
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Panic().Err(err).Msg("Unable to start web server")
		os.Exit(1)
	}
}
