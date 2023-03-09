package blog

import (
	"kurisu/web"
)

type blog struct {
	Title string
}

var b = blog{}

func New(msgBus chan string) {
	web.RegisterRoute(&b)
}
