package main

import (
	"github.com/jackc/pgx"
	"log"
)

func mustPrepare(db *pgx.Conn, name, query string) *pgx.PreparedStatement {
	stmt, err := db.Prepare(name, query)
	if err != nil {
		log.Fatalf("Error when preparing statement %q: %s", query, err)
	}
	return stmt
}

// design and code by tsingson
