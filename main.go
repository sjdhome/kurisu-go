package main

import (
	"net/http"
	"log"
	"flag"
	"os"
)

type BlogPostHandler struct {}

func (h BlogPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
	log.Println("Receive request")
}

var logFile string

func main() {
	flag.StringVar(&logFile, "log-file", "", "Specify log storage directory.")
	flag.Parse()
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Unable to open log file.")
			log.Fatalln(err)
		}
		log.SetOutput(f)
	}

	var blogPostHandler BlogPostHandler
	http.Handle("/blog/post/", blogPostHandler)
	log.Println("HTTP server started at :3000.")
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
