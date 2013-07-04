package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
)

type InboundEmailHandler struct{}

func (h *InboundEmailHandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Print("Could not parse email: " + err.Error())
		http.Error(w, "Could not parse email: "+err.Error(), 500)
		return
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Print("Could not read body: " + err.Error())
		http.Error(w, "Could not read body: "+err.Error(), 500)
		return
	}
	// Body is an actual email
	log.Print(msg.Header.Get("From"))
	log.Print(msg.Header.Get("Subject"))
	log.Print(string(body))
}

func main() {
	addr := os.Getenv("BOREDGAME_EMAIL_SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:9999"
	}
	log.Print("Running incoming email server on http://" + addr)
	err := http.ListenAndServe(addr, &InboundEmailHandler{})
	if err != nil {
		panic(err.Error())
	}
}
