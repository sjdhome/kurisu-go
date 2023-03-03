package terminal

import (
	"fmt"
	"log"
)

func New(msgBus chan string) {
	var input string
	for {
		fmt.Scan(&input)
		switch input {
		case "exit":
			log.Println("Goodbye!")
			msgBus <- "exit"
		case "help":
			log.Println("Available commands:")
			log.Printf("\texit: Exit the program.\n")
			log.Printf("\thelp: Show this help message.\n")
		default:
			log.Printf("Unknown command \"%s\", ignoring.", input)
		}
	}
}
