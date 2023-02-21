package blog

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

type HTTPHandler struct {
}

//go:embed page
var page embed.FS

//go:embed static
var static embed.FS
var staticHandler = http.FileServer(http.FS(static))

func (handler HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	log.Printf("request %s\n", p)
	switch {
	case strings.HasPrefix(p, "/static/"):
		staticHandler.ServeHTTP(w, r)
	case strings.HasPrefix(p, "/"):
		b, err := page.ReadFile("page/index.html")
		if err != nil {
			log.Panicln(err)
		}
		w.Write(b)
	}
}
