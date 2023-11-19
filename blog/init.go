package blog

import (
	"net/http"
)

func Init() {
	http.Handle("/blog/", BlogHandler{
		blog: &Blog{},
	})
}
