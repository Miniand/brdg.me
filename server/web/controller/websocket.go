package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/Miniand/brdg.me/server/communicate"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to websocket, %s", err)
		return
	}
	token := ""
	if err := conn.ReadJSON(&token); err != nil {
		conn.WriteMessage(
			websocket.TextMessage,
			websocket.FormatCloseMessage(
				websocket.CloseUnsupportedData,
				"We could not understand your authorization token, please send a plain JSON string with your token after connecting",
			),
		)
		conn.Close()
		return
	}
	communicate.RegisterPlayerConnection(token, conn)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
	communicate.UnregisterPlayerConnection(token, conn)
}
