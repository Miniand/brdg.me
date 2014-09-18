package main

import (
	"github.com/Miniand/brdg.me/server/email/parser"
	"github.com/Miniand/brdg.me/server/web"
)

func main() {
	result := make(chan error)
	go func() {
		result <- web.Run()
	}()
	go func() {
		result <- parser.Run()
	}()
	if err := <-result; err != nil {
		panic(err.Error())
	}
}
