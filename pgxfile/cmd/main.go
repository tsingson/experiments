package main

import (
	"fmt"
	"os"

	"github.com/tsingson/experiments/pgxfile"
	"github.com/tsingson/uuid"
)

func main() {
	conn := utils.Connect("connect test")
	defer conn.Close()
	fmt.Printf("Connection worked!\n")

	_, err := conn.Exec(`
	  	create table users (
	      id uuid constraint user_pk primary key   not null,
	      username text constraint unique_username unique not null,
	      password text not null,  -- wow, this is totally not secure
	      first_name text not null,
	      last_name text not null);
	  	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create users table: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully created users table\n")
	fmt.Println(uuid.NewV4().String())

	id := uuid.NewV4()

	_, err = conn.Exec(`
	insert into users (
    id,
    username,
    password,
    first_name,
    last_name) values ($1, $2, $3, $4, $5);
	`, id, "mwood", "passwd", "Manni", "Wood")
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "unique_username" {
				fmt.Fprintf(os.Stderr, "Unable to create user mwood, because username already taken: %v\n", pgerr)
			} else {
				fmt.Fprintf(os.Stderr, "Unexpected postgres error trying to create user mwood: %v\n", pgerr)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unexpected error trying to create user mwood: %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Printf("Successfully created user mwood\n")
	fmt.Printf("Successfully created user mwood\n")
}
