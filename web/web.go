package web

import (
	"embed"
	"flag"
	"io/fs"
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

//go:embed www
var _wwwFS embed.FS

var wwwFS, _ = fs.Sub(_wwwFS, "www")
var wwwServer = http.FileServer(http.FS(wwwFS))

//go:embed favicon.ico
var favicon []byte

//go:embed robots.txt
var robots []byte

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
		http.Redirect(w, r, r.Proto+"://"+DOMAIN+r.URL.Path, http.StatusMovedPermanently)
		return
	}

	// Router.
	switch {
	case p == "/favicon.ico":
		w.WriteHeader(http.StatusOK)
		h.Set("Content-Type", "image/x-icon")
		w.Write(favicon)
	case p == "/robots.txt":
		w.WriteHeader(http.StatusOK)
		h.Set("Content-Type", "text/plain")
		w.Write(robots)
	default:
		h.Set("Access-Control-Allow-Origin", r.Proto+"://"+r.Host)
		if strings.HasSuffix(p, "/") {
			p += "index.html"
		}
		f, err := wwwFS.Open(p[1:])
		if err != nil {
			log.Println("\t", err)
			http.NotFound(w, r)
			return
		}
		s, err := f.Stat()
		if err != nil {
			log.Println("\t", err)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		if s.IsDir() {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
			return
		}
		wwwServer.ServeHTTP(w, r)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
