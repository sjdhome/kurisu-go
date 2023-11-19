package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var port int
var dbPath string

func main() {
	flag.IntVar(&port, "port", 3000, "HTTP port")
	flag.StringVar(&dbPath, "db-path", "blog.sqlite3", "SQLite database path")
	flag.Parse()

	db, err := sql.Open("sqlite3", dbPath)
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
	log.Printf("HTTP server started at :%d.\n", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
