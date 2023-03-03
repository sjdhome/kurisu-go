package web

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

func New(msgBus chan string) {
	log.Println("Starting web server at \":8080\".")
	http.ListenAndServe(":8080", http.HandlerFunc(serve))
}

//go:embed static
var staticFS embed.FS

var staticServer = http.FileServer(http.FS(staticFS))

//go:embed page/index.html
var indexHTML []byte

func serve(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
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
