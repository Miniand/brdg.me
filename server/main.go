package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		// Default to port 9999
		addr = "localhost:9999"
	}
	http.HandleFunc("/inbound", InboundEmailHandler)
	http.HandleFunc("/", RootHandler)
	fmt.Println("Running server on http://" + addr)
	http.ListenAndServe(addr, nil)
}

func InboundEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'll eventually handle email responses")
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "LALALALA I'm A web server!")
}
