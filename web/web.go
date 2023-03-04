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

//go:embed favicon.ico
var favicon []byte

const DOMAIN = "sjdhome.com"

func serve(w http.ResponseWriter, r *http.Request) {
	if ip, exist := r.Header[http.CanonicalHeaderKey("CF-Connecting-IP")]; exist {
		r.RemoteAddr = ip[0]
	}
	log.Printf("%s: %s %s %s\n", r.Host, r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method == http.MethodGet {
		get(w, r)
	} else if r.Method == http.MethodPost {
		post(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	h := w.Header()

	h.Set("Server", "kurisu")

	// Redirect www to naked domain.
	if r.Host == "www."+DOMAIN {
		http.Redirect(w, r, "https://"+DOMAIN+r.URL.Path, http.StatusMovedPermanently)
		return
	}

	// Router.
	switch {
	case p == "/" || p == "/index.html":
		w.WriteHeader(http.StatusOK)
		h.Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	case strings.HasPrefix(p, "/static/"):
		h.Set("Access-Control-Allow-Origin", "https://"+DOMAIN)
		staticServer.ServeHTTP(w, r)
	case p == "/favicon.ico":
		w.WriteHeader(http.StatusOK)
		h.Set("Content-Type", "image/x-icon")
		w.Write(favicon)
	default:
		http.NotFound(w, r)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
