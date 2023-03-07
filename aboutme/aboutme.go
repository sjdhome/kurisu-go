package aboutme

import (
	"embed"
	"io/fs"
	"kurisu/web"
	"net/http"
	"strings"
)

//go:embed www
var _wwwFS embed.FS
var wwwFS, _ = fs.Sub(_wwwFS, "www")
var wwwServer = http.FileServer(http.FS(wwwFS))

type aboutMe struct{}

func (a aboutMe) Method() string {
	return "GET"
}

func (a aboutMe) Path(path string) bool {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	f, err := wwwFS.Open(path)
	if err == nil {
		s, err := f.Stat()
		if err == nil && s.IsDir() {
			return false
		}
	}
	if strings.HasSuffix(path, ".ts") {
		return false
	}
	return path == "" ||
		path == "index.html" ||
		strings.HasPrefix(path, "js/") ||
		strings.HasPrefix(path, "css/")
}

func (a aboutMe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wwwServer.ServeHTTP(w, r)
}

var route = aboutMe{}

func New(msgBus chan string) {
	web.RegisterRoute(&route)
}
