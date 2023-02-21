package main

import (
	"kurisu/app/blog"
	"kurisu/app/terminal"
)

func main() {
	ch := make(chan string)
	var msg string
	go blog.Run()
	go terminal.Run(ch)
	for ; msg != "exit"; msg = <-ch {
		// cope with msg
	}
}
