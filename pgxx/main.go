package main

import (
	"github.com/jackc/pgx"
	"log"
)

var conn *pgx.Conn

func main() {
	/**
	     config, err := pgx.ParseEnvLibpq()

		if err != nil {
			log.Fatalln("Unable to extract config from env:", err)
		}
	*/
	config := pgx.ConnConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Database: "testdb2",
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	var projectID int64
	err = conn.QueryRow("SELECT id FROM shard_1.projects WHERE user_id = $1", int64(1341396671631721500)).Scan(&projectID)
	if err != nil {
		log.Fatalln("Query failed:", err)
	}

	log.Println("Succeeded:", projectID)
}
