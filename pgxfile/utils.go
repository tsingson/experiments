package utils

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func Connect(applicationName string) (conn *pgx.Conn) {
	var runtimeParams map[string]string
	runtimeParams = make(map[string]string)
	runtimeParams["application_name"] = applicationName
	connConfig := pgx.ConnConfig{
		User:              "postgres",
		Password:          "postgres",
		Host:              "localhost",
		Port:              5432,
		Database:          "testdb",
		TLSConfig:         nil,
		UseFallbackTLS:    false,
		FallbackTLSConfig: nil,
		RuntimeParams:     runtimeParams,
	}
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}
	return conn
}
