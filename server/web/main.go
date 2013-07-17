package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", RootHandler)
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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to brdg.me!  Send an email to play@brdg.me to start playing board games over email!")
}
