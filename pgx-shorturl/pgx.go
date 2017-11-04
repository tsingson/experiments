package main

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/log15adapter"
	log "gopkg.in/inconshreveable/log15.v2"
)

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("getUrl", `
    select url from shortened_urls where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("deleteUrl", `
    delete from shortened_urls where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("putUrl", `
    insert into shortened_urls(id, url) values ($1, $2)
    on conflict (id) do update set url=excluded.url
  `)
	return
}

func fetchUrl(path string) (string, error) {
	var url string
	err := pool.QueryRow("getUrl", path).Scan(&url)
	return url, err
}

func saveUrl(id, url string) (pgx.CommandTag, error) {
	return pool.Exec("putUrl", id, url)
}
func deleteUrl(path string) (pgx.CommandTag, error) {
	return pool.Exec("deleteUrl", path)
}

func initPgx() (*pgx.ConnPool, error) {
	logger := log15adapter.NewLogger(log.New("module", "pgx"))
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "127.0.0.1",
			User:     "postgres",
			Password: "postgres",
			Database: "test",
			Logger:   logger,
		},
		MaxConnections: 5,
		AfterConnect:   afterConnect,
	}
	return pgx.NewConnPool(connPoolConfig)

}
