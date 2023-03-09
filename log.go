package main

import (
	"io"
	"log"
	"os"
)

func initLog() {
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
}
