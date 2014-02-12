package main

import (
	"github.com/Miniand/brdg.me/server/email"
	"github.com/Miniand/brdg.me/server/web"
)

func main() {
	result := make(chan error)
	go func() {
		result <- web.Run()
	}()
	go func() {
		result <- email.Run()
	}()
	if err := <-result; err != nil {
		panic(err.Error())
	}
}
