package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
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
	player := ParseFrom(msg.Header.Get("From"))
	log.Print("Player:", player)
	gameId := ParseSubject(msg.Header.Get("Subject"))
	log.Print("Game ID:", gameId)
	commands := ParseBody(string(body))
	log.Print("Commands:", strings.Join(commands, ", "))
	err = HandleCommands(player, gameId, commands)
	if err != nil {
		log.Print("Error handling commands: " + err.Error())
		http.Error(w, "Error handling commands: "+err.Error(), 500)
		return
	}
}

func main() {
	addr := os.Getenv("BOREDGAME_EMAIL_SERVER_ADDRESS")
	if addr == "" {
		addr = ":9999"
	}
	log.Print("Running incoming email server on " + addr)
	err := http.ListenAndServe(addr, &InboundEmailHandler{})
	if err != nil {
		panic(err.Error())
	}
}
