package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sjdhome.com/kurisu-go/blog"
	"sjdhome.com/kurisu-go/database"
)

var port int

func main() {
	flag.IntVar(&port, "port", 3000, "HTTP server listen port")
	flag.Parse()

	err := database.Init()
	if err != nil {
		panic(err)
	}
	blog.Init()

	log.Printf("HTTP server started at :%d.\n", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
