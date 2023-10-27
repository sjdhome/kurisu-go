package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

var logFile string
var port int
var dbFile string

func main() {
	flag.StringVar(&logFile, "log-file", "", "Specify log storage directory.")
	flag.IntVar(&port, "port", 3000, "HTTP server port.")
	flag.StringVar(&dbFile, "db-file", "kurisu.sqlite3", "SQLite database file.")
	flag.Parse()

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("Unable to open log file: %w", err))
		}
		log.SetOutput(f)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(fmt.Errorf("Unable to open database: %w", err))
	}
	defer db.Close()

	blog := Blog{
		db: db,
	}
	blogHandler := BlogHandler{
		blog: &blog,
	}
	http.Handle("/blog/", blogHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	log.Printf("HTTP server started at :%d.\n", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
