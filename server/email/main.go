package main

import (
	"io/ioutil"
	"labix.org/v2/mgo/bson"
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
	subject := msg.Header.Get("Subject")
	if bson.IsObjectIdHex(subject) {
		// Load a game and play it
		gm, err := LoadGame(bson.ObjectIdHex(subject))
		if err != nil {
			log.Print("Could not load game: " + err.Error())
			http.Error(w, "Could not load game: "+err.Error(), 500)
		}

		_, err = gm.ToGame()
		if err != nil {
			log.Print("Could not read game: " + err.Error())
			http.Error(w, "Could not read game: "+err.Error(), 500)
		}
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
