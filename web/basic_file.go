package web

import (
	_ "embed"
	"net/http"
)

//go:embed favicon.ico
var favicon []byte

//go:embed robots.txt
var robots []byte

//go:embed sitemap.txt
var sitemap []byte

type basicFile struct{}

func (f *basicFile) Method() string {
	return "GET"
}

func (f *basicFile) Path(path string) bool {
	return path == "/favicon.ico" || path == "/robots.txt" || path == "/sitemap.txt"
}

func (f *basicFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/favicon.ico":
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(favicon)
	case "/robots.txt":
		w.Header().Set("Content-Type", "text/plain")
		w.Write(robots)
	case "/sitemap.txt":
		w.Header().Set("Content-Type", "text/plain")
		w.Write(sitemap)
	}
}

func registerBasicRoutes() {
	RegisterRoute(&basicFile{})
}
