package blog

import (
	"embed"
	"html/template"
	"io/fs"
	"kurisu/web"
	"log"
	"net/http"
	"strings"
)

type blog struct{}

//go:embed www
var _wwwFS embed.FS
var wwwFS, _ = fs.Sub(_wwwFS, "www")
var wwwServer = http.StripPrefix("/blog", http.FileServer(http.FS(wwwFS)))

func (b blog) Method() string {
	return "GET"
}

func (b blog) Path(path string) bool {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	f, err := wwwFS.Open(path)
	if err == nil {
		s, err := f.Stat()
		if err == nil && s.IsDir() {
			return false
		}
	}
	return strings.HasPrefix(path, "blog")
}

func (b blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimSuffix(p, "/")
	if p == "blog" || p == "blog/index.html" {
		tmpl, err := template.ParseFS(wwwFS, "index.html")
		if err != nil {
			log.Println("\t", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, b)
		if err != nil {
			log.Println("\t", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	wwwServer.ServeHTTP(w, r)
}

var route = blog{}

func New(msgBus chan string) {
	web.RegisterRoute(&route)
}
