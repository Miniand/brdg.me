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
		log.Println("Could not parse email: " + err.Error())
		http.Error(w, "Could not parse email: "+err.Error(), 500)
		return
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Println("Could not read body: " + err.Error())
		http.Error(w, "Could not read body: "+err.Error(), 500)
		return
	}
	// Body is an actual email
	player := ParseFrom(msg.Header.Get("From"))
	log.Println("Player:", player)
	gameId := ParseSubject(msg.Header.Get("Subject"))
	log.Println("Game ID:", gameId)
	commands := ParseBody(string(body))
	log.Println("Commands:")
	for _, c := range commands {
		log.Println("Commands:", strings.Join(c, " "))
	}
	err = HandleCommands(player, gameId, commands)
	if err != nil {
		log.Println("Error handling commands: " + err.Error())
		http.Error(w, "Error handling commands: "+err.Error(), 500)
		return
	}
}

func main() {
	addr := os.Getenv("BOREDGAME_EMAIL_SERVER_ADDRESS")
	if addr == "" {
		addr = ":9999"
	}
	log.Println("Running incoming email server on " + addr)
	err := http.ListenAndServe(addr, &InboundEmailHandler{})
	if err != nil {
		panic(err.Error())
	}
}
