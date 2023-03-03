package web

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"
)

func New(msgBus chan string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting web server at \":%s\".\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, http.HandlerFunc(serve)))
}

//go:embed static
var staticFS embed.FS

var staticServer = http.FileServer(http.FS(staticFS))

//go:embed page/index.html
var indexHTML []byte

func serve(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	w.Header().Set("Server", "kurisu")
	switch {
	case p == "/" || p == "/index.html":
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	case strings.HasPrefix(p, "/static/"):
		staticServer.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
