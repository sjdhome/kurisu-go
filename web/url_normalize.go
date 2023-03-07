package web

import (
	"io/fs"
	"log"
	"net/http"
	"strings"
)

func URLNormalize(w http.ResponseWriter, r *http.Request, prefix string, fs fs.FS) bool {
	p, found := strings.CutPrefix(r.URL.Path, prefix)
	if !found {
		log.Println("URLNormalize: prefix not found", r.URL.Path, prefix)
		http.NotFound(w, r)
		return false
	}
	if strings.HasSuffix(p, "/") || p == "" {
		log.Println("URLNormalize: redirect to index.html", r.URL.Path)
		p += "index.html"
		r.URL.Path += "index.html"
	}
	f, err := fs.Open(p)
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println("URLNormalize: failed to close file", p)
			log.Println("\t", err)
		}
	}()
	if err != nil {
		log.Println("URLNormalize: failed to open file", p)
		http.NotFound(w, r)
		return false
	}
	s, err := f.Stat()
	if err != nil {
		log.Println("URLNormalize: failed to stat file", p)
		http.NotFound(w, r)
		return false
	}
	if s.IsDir() {
		log.Println("URLNormalize: redirect to directory", r.URL.Path)
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return false
	} else if strings.HasSuffix(s.Name(), ".ts") {
		log.Println("URLNormalize: user is trying to acess typescript file, ignoring.", r.URL.Path)
		http.NotFound(w, r)
		return false
	}
	return true
}
