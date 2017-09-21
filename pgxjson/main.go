package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/json-iterator/go"
	"log"
)

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    contact JSON
);`

type Person struct {
	FirstName string            `db:"first_name"`
	LastName  string            `db:"last_name"`
	Contact   map[string]string `db:"contact"`
}

func main() {
	config := pgx.ConnConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Database: "test",
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// exec the schema or fail
	_, err = conn.Exec("DROP TABLE IF EXISTS person;")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(schema)
	if err != nil {
		panic(err)
	}

	var person Person
	contact := map[string]string{"email": "jmoiron@jmoiron.net"}
	contact_json, _ := jsoniter.Marshal(contact)

	_, err = conn.Exec("INSERT INTO person (first_name, last_name, contact) VALUES ($1, $2, $3)", "Jason", "Moiron", contact_json)
	if err != nil {
		panic(err)
	}

	row := conn.QueryRow("SELECT first_name, last_name, contact FROM person LIMIT 1")
	err = row.Scan(&person.FirstName, &person.LastName, &person.Contact)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", person)
}
