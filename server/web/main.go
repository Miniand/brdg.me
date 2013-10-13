package main

import (
	"github.com/Miniand/brdg.me/server/web/route"
	"log"
	"net/http"
	"os"
)

func main() {
	http.Handle("/", route.Router())
	addr := os.Getenv("BRDGME_WEB_SERVER_ADDRESS")
	if addr == "" {
		addr = ":9998"
	}
	log.Print("Running web server on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
