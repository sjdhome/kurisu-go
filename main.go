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
var blogDB string

func main() {
	flag.StringVar(&logFile, "log-file", "", "Specify log storage directory.")
	flag.IntVar(&port, "port", 3000, "HTTP server port.")
	flag.StringVar(&blogDB, "blogDB", "blog.sqlite3", "Blog SQLite database file.")
	flag.Parse()

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("Unable to open log file: %w", err))
		}
		defer f.Close()
		log.SetOutput(f)
	}

	db, err := sql.Open("sqlite3", blogDB)
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
