package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", RootHandler)
	addr := os.Getenv("BOREDGAME_WEB_SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:9998"
	}
	log.Print("Running web server on http://" + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "LALALALA I'm A web server!")
}
