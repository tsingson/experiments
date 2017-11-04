package main

import (
	"flag"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/tsingson/experiments/fastpgx/templates"
	"github.com/valyala/fasthttp"
	"log"
	"runtime"
	"sort"
)

var (
	worldSelectStmt   *pgx.PreparedStatement
	worldUpdateStmt   *pgx.PreparedStatement
	fortuneSelectStmt *pgx.PreparedStatement

	db *pgx.ConnPool
)

func main() {
	flag.Parse()

	var err error
	maxConnectionCount := runtime.NumCPU() * 2
	if db, err = initDatabase("localhost", "benchmarkdbuser", "benchmarkdbpass", "hello_world", 5432, maxConnectionCount); err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	s := &fasthttp.Server{
		Handler: mainHandler,
		Name:    "go",
	}
	ln := GetListener()
	if err = s.Serve(ln); err != nil {
		log.Fatalf("Error when serving incoming connections: %s", err)
	}
}

func mainHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	switch string(path) {
	case "/plaintext":
		PlaintextHandler(ctx)
	case "/json":
		JSONHandler(ctx)
	case "/db":
		dbHandler(ctx)
	case "/queries":
		queriesHandler(ctx)
	case "/fortune":
		fortuneHandler(ctx)
	case "/update":
		updateHandler(ctx)
	default:
		ctx.Error("unexpected path", fasthttp.StatusBadRequest)
	}
}

func dbHandler(ctx *fasthttp.RequestCtx) {
	var w World
	fetchRandomWorld(&w)
	JSONMarshal(ctx, &w)
}

func queriesHandler(ctx *fasthttp.RequestCtx) {
	n := GetQueriesCount(ctx)
	worlds := make([]World, n)
	for i := 0; i < n; i++ {
		fetchRandomWorld(&worlds[i])
	}
	JSONMarshal(ctx, worlds)
}

func fortuneHandler(ctx *fasthttp.RequestCtx) {
	rows, err := db.Query("fortuneSelectStmt")
	if err != nil {
		log.Fatalf("Error selecting db data: %v", err)
	}

	fortunes := make([]templates.Fortune, 0, 16)
	for rows.Next() {
		var f templates.Fortune
		if err := rows.Scan(&f.ID, &f.Message); err != nil {
			log.Fatalf("Error scanning fortune row: %s", err)
		}
		fortunes = append(fortunes, f)
	}
	rows.Close()
	fortunes = append(fortunes, templates.Fortune{Message: "Additional fortune added at request time."})

	sort.Sort(FortunesByMessage(fortunes))

	ctx.SetContentType("text/html; charset=utf-8")
	templates.WriteFortunePage(ctx, fortunes)
}

func updateHandler(ctx *fasthttp.RequestCtx) {
	n := GetQueriesCount(ctx)

	worlds := make([]World, n)
	for i := 0; i < n; i++ {
		w := &worlds[i]
		fetchRandomWorld(w)
		w.RandomNumber = int32(RandomWorldNum())
	}

	// sorting is required for insert deadlock prevention.
	sort.Sort(WorldsByID(worlds))
	txn, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %s", err)
	}

	for i := 0; i < n; i++ {
		w := &worlds[i]
		if _, err = txn.Exec("worldUpdateStmt", w.RandomNumber, w.Id); err != nil {
			log.Fatalf("Error updating world row %d: %s", i, err)
		}
	}
	if err = txn.Commit(); err != nil {
		log.Fatalf("Error when commiting world rows: %s", err)
	}

	JSONMarshal(ctx, worlds)
}

func fetchRandomWorld(w *World) {
	n := RandomWorldNum()
	if err := db.QueryRow("worldSelectStmt", n).Scan(&w.Id, &w.RandomNumber); err != nil {
		log.Fatalf("Error scanning world row: %s", err)
	}
}

func initDatabase(dbHost string, dbUser string, dbPass string, dbName string, dbPort uint16, maxConnectionsInPool int) (*pgx.ConnPool, error) {

	var successOrFailure string = "OK"
	var config pgx.ConnPoolConfig

	config.Host = dbHost
	config.User = dbUser
	config.Password = dbPass
	config.Database = dbName
	config.Port = dbPort

	config.MaxConnections = maxConnectionsInPool

	config.AfterConnect = func(conn *pgx.Conn) error {
		worldSelectStmt = mustPrepare(conn, "worldSelectStmt", "SELECT id, randomNumber FROM World WHERE id = $1")
		worldUpdateStmt = mustPrepare(conn, "worldUpdateStmt", "UPDATE World SET randomNumber = $1 WHERE id = $2")
		fortuneSelectStmt = mustPrepare(conn, "fortuneSelectStmt", "SELECT id, message FROM Fortune")
		return nil
	}

	fmt.Println("--------------------------------------------------------------------------------------------")

	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		successOrFailure = "FAILED"
		log.Println("Connecting to database ", dbName, " as user ", dbUser, " ", successOrFailure, ": \n ", err)
	} else {
		log.Println("Connecting to database ", dbName, " as user ", dbUser, ": ", successOrFailure)

		log.Println("Fetching one record to test if db connection is valid...")
		var w World
		n := RandomWorldNum()
		if errPing := connPool.QueryRow("worldSelectStmt", n).Scan(&w.Id, &w.RandomNumber); errPing != nil {
			log.Fatalf("Error scanning world row: %s", errPing)
		}
		log.Println("OK")
	}

	fmt.Println("--------------------------------------------------------------------------------------------")

	return connPool, err
}
