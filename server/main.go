package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
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
	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		fmt.Println("Could not parse email: " + err.Error())
		http.Error(w, "Could not parse email: "+err.Error(), 500)
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		fmt.Println("Could not read body: " + err.Error())
		http.Error(w, "Could not read body: "+err.Error(), 500)
	}
	// Body is an actual email
	fmt.Println(msg.Header.Get("From"))
	fmt.Println(msg.Header.Get("Subject"))
	fmt.Println(string(body))
	fmt.Fprintf(w, msg.Header.Get("From"))
	fmt.Fprintf(w, msg.Header.Get("Subject"))
	fmt.Fprintf(w, string(body))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "LALALALA I'm A web server!")
}
