package terminal

import "fmt"

func Run(ch chan string) {
	var s string
	for {
		fmt.Scanln(&s)
		switch s {
		case "exit":
			ch <- s
		}
	}
}
