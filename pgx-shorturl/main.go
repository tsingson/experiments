package main

import (
	"github.com/jackc/pgx"

	log "gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	pool *pgx.ConnPool
	err  error
)

func getUrlHandler(w http.ResponseWriter, req *http.Request) {
	url, err := fetchUrl(req.URL.Path)
	switch err {
	case nil:
		http.Redirect(w, req, url, http.StatusSeeOther)
	case pgx.ErrNoRows:
		http.NotFound(w, req)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func putUrlHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path
	var url string
	if body, err := ioutil.ReadAll(req.Body); err == nil {
		url = string(body)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := saveUrl(id, url); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func deleteUrlHandler(w http.ResponseWriter, req *http.Request) {
	if _, err := deleteUrl(req.URL.Path); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func urlHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		getUrlHandler(w, req)

	case "PUT":
		putUrlHandler(w, req)

	case "DELETE":
		deleteUrlHandler(w, req)

	default:
		w.Header().Add("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {

	pool, err = initPgx()
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	http.HandleFunc("/", urlHandler)

	log.Info("Starting URL shortener on localhost:8080")
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Crit("Unable to start web server", "error", err)
		os.Exit(1)
	}
}
