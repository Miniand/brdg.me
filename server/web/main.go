package main

import (
	"github.com/Miniand/brdg.me/server/web/config"
	"github.com/Miniand/brdg.me/server/web/route"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", route.Router())
	addr := config.Get(config.SERVER_ADDRESS)
	log.Print("Running web server on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
