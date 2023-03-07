package web

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
	"strings"
)

var (
	webAddr = flag.String("web", ":8080", "Listen address.")
	cert    = flag.String("cert", "", "TLS certificate file.")
	key     = flag.String("key", "", "TLS key file.")
)

func New(msgBus chan string) {
	if *webAddr == "" {
		return
	}

	registerBasicRoutes()

	if *cert != "" && *key != "" {
		log.Printf("Starting HTTPS server at \"%s\".\n", *webAddr)
		log.Fatalln(http.ListenAndServeTLS(*webAddr, *cert, *key, http.HandlerFunc(serve)))
	} else {
		log.Printf("Starting HTTP server at \"%s\".\n", *webAddr)
		log.Fatalln(http.ListenAndServe(*webAddr, http.HandlerFunc(serve)))
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	// Get real IP address.
	if ip, exist := r.Header[http.CanonicalHeaderKey("CF-Connecting-IP")]; exist {
		r.RemoteAddr = ip[0]
	}

	log.Printf("%s: %s %s %s\n", r.Host, r.RemoteAddr, r.Method, r.URL.Path)

	r.Header["Server"] = []string{"kurisu"}

	// Redirect www to naked domain.
	if strings.HasPrefix(r.Host, "www.") {
		naked := strings.TrimPrefix(r.Host, "www.")
		http.Redirect(w, r, r.Proto+"://"+naked+r.URL.Path, http.StatusMovedPermanently)
		return
	}

	route, err := SelectRoute(r.URL.Path)
	if err != nil {
		log.Println("\t", err)
		http.NotFound(w, r)
		return
	}
	route.ServeHTTP(w, r)
}
