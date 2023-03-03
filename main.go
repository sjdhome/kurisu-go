package main

import (
	"fmt"
	"kurisu/terminal"
	"kurisu/web"
	"log"
	"os"
)

func main() {
	log.Println("Starting kurisu...")

	webMsg := make(chan string)
	go web.New(webMsg)

	terminalMsg := make(chan string)
	go terminal.New(terminalMsg)

	log.Println("Now you can type command. Type \"help\" to get help.")
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
