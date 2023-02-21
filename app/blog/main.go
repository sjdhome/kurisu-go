package blog

import (
	"net/http"
)

func Run() {
	var h HTTPHandler
	http.ListenAndServe(":8080", h)
}
