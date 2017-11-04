package main

import (
	"encoding/hex"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"log"
)

var schema = `
CREATE TABLE person (
     id uuid constraint user_pk primary key  not null,
    first_name text,
    last_name text,
    contact JSON
);`

type (
	Person struct {
		Id        string  `db:"id"`
		FirstName string  `db:"first_name"`
		LastName  string  `db:"last_name"`
		Contact   Contact `db:"contact"`
	}
	Contact struct {
		Address string `json:"address"`
		Post    string `json:"post"`
	}
)

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

	id := NewGuid()
	var person Person
	contact := new(Contact)
	contact.Address = "tsingson"
	contact.Post = "1234"
	contact_json, _ := jsoniter.Marshal(contact)

	_, err = conn.Exec("INSERT INTO person (id, first_name, last_name, contact) VALUES ($1, $2, $3, $4)", id, "Jason", "Moiron", contact_json)
	if err != nil {
		panic(err)
	}

	row := conn.QueryRow("SELECT id, first_name, last_name, contact::json FROM person LIMIT 1")
	err = row.Scan(&person.Id, &person.FirstName, &person.LastName, &person.Contact)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", person.Contact.Post)
	fmt.Printf("%v", person.Id)
}

type Guid struct {
	Uuid uuid.UUID
}

func Uuid4String() string {
	return NewGuid()
}

//
func NewGuid() string {
	//	enc := new(guid.Base58)
	u := uuid.NewV4()
	return hex.EncodeToString(u.Bytes())

	//return base58.Encode([]byte(ustring))
}
