package main

import (
	"flag"
	"fmt"
	"kurisu/aboutme"
	"kurisu/blog"
	"kurisu/terminal"
	"kurisu/web"
	"os"
)

var logFilename = flag.String("log", "kurisu.log", "Log filename.")

func main() {
	flag.Parse()

	initLog()

	webMsg := make(chan string)
	go web.New(webMsg)

	terminalMsg := make(chan string)
	go terminal.New(terminalMsg)

	blogMsg := make(chan string)
	go blog.New(blogMsg)

	aboutmeMsg := make(chan string)
	go aboutme.New(aboutmeMsg)

	select {
	case msg := <-webMsg:
		fmt.Println(msg)
	case msg := <-terminalMsg:
		switch msg {
		case "exit":
			os.Exit(0)
		}
	}
}
