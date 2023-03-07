package web

import (
	"errors"
	"net/http"
)

type Route interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Method() string
	Path(string) bool
}

var routeTable = []Route{}

func RegisterRoute(r Route) {
	routeTable = append(routeTable, r)
}

func SelectRoute(path string) (Route, error) {
	for _, r := range routeTable {
		if r.Path(path) {
			return r, nil
		}
	}
	return nil, errors.New("Route not found")
}
