package web

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"strings"
)

var (
	listenAt = flag.String("l", ":8080", "Listen address.")
	certfile = flag.String("cert", "", "TLS certificate file.")
	keyfile  = flag.String("key", "", "TLS key file.")
)

func New(msgBus chan string) {
	if *certfile != "" && *keyfile != "" {
		log.Printf("Starting HTTPS server at \"%s\".\n", *listenAt)
		log.Fatalln(http.ListenAndServeTLS(*listenAt, *certfile, *keyfile, http.HandlerFunc(serve)))
	} else {
		log.Printf("Starting HTTP server at \"%s\".\n", *listenAt)
		log.Fatalln(http.ListenAndServe(*listenAt, http.HandlerFunc(serve)))
	}
}

//go:embed static
var staticFS embed.FS

var staticServer = http.FileServer(http.FS(staticFS))

//go:embed page/index.html
var indexHTML []byte

func serve(w http.ResponseWriter, r *http.Request) {
	log.Printf("web: %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

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
