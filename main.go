package main

import (
	"flag"
	"fmt"
	"io"
	"kurisu/terminal"
	"kurisu/web"
	"log"
	"os"
)

var logFilename = flag.String("log", "kurisu.log", "Log filename.")

func main() {
	flag.Parse()

	// Send log to stdout and file.
	logFile, err := os.OpenFile(*logFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open log file.")
		log.Fatalln(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.SetPrefix("[kurisu] ")
	log.Println("Starting kurisu...")

	webMsg := make(chan string)
	go web.New(webMsg)

	terminalMsg := make(chan string)
	go terminal.New(terminalMsg)

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
