package main

import (
	"errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/zerologadapter"
	"github.com/rs/zerolog"
)

// AfterConnect creates the prepared statements that this application uses

var (
	prepare map[string]string
)

func init() {
	prepare = map[string]string{
		"getUrl":    `select url from shortened_urls where id=$1`,
		"deleteUrl": `delete from shortened_urls where id=$1`,
		"putUrl":    `insert into shortened_urls(id, url) values ($1, $2) on conflict (id) do update set url=excluded.url`}

}

func afterConnectMap(conn *pgx.Conn) error {

	if len(prepare) == 0 {
		err = errors.New("null map ")
		return err
	}
	for key, value := range prepare {
		_, err = conn.Prepare(key, value)
		if err != nil {
			return err
		}
	}
	return nil

}

func AfterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("getUrl", `
    select url from shortened_urls where id=$1
  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare("deleteUrl", `
    delete from shortened_urls where id=$1
  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare("putUrl", `
    insert into shortened_urls(id, url) values ($1, $2)
    on conflict (id) do update set url=excluded.url
  `)
	if err != nil {
		return err
	}
	return nil
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

func InitPgx(log zerolog.Logger) (*pgx.ConnPool, error) {
	logger := zerologadapter.NewLogger(log)
	pgxConfig := pgx.ConnConfig{
		Host:     "127.0.0.1",
		User:     "postgres",
		Password: "postgres",
		Database: "test",
		Logger:   logger,
	}

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgxConfig,
		MaxConnections: 5,
		AfterConnect:   afterConnectMap,
	}
	return pgx.NewConnPool(connPoolConfig)
}
