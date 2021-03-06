package parser

import (
	"log"
	"net/http"
	"os"

	"github.com/Miniand/brdg.me/server/game"
)

type InboundEmailHandler struct{}

func (h *InboundEmailHandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	msg, body, err := GetPlainEmailBody(r.Body)
	if err != nil {
		log.Println("Could not parse email: " + err.Error())
		http.Error(w, "Could not parse email: "+err.Error(), 500)
		return
	}
	// Body is an actual email
	player := ParseFrom(msg.Header.Get("From"))
	gameId := ParseSubject(msg.Header.Get("Subject"))
	commandText := ParseBody(body)
	err = game.HandleCommandText(player, gameId, commandText)
	if err != nil {
		log.Println("Error handling commands: " + err.Error())
		http.Error(w, "Error handling commands: "+err.Error(), 500)
		return
	}
}

func Run() error {
	addr := os.Getenv("BRDGME_EMAIL_SERVER_ADDRESS")
	if addr == "" {
		addr = ":9999"
	}
	log.Println("Running incoming email server on " + addr)
	return http.ListenAndServe(addr, &InboundEmailHandler{})
}
