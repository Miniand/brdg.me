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
	user, found, err := AuthenticateToken(token)
	if err != nil {
		log.Printf("Error authenticating token over websocket, %s", err)
		conn.WriteMessage(
			websocket.TextMessage,
			websocket.FormatCloseMessage(
				websocket.CloseUnsupportedData,
				"There was an error authenticating the token",
			),
		)
		conn.Close()
		return
	}
	if !found {
		conn.WriteMessage(
			websocket.TextMessage,
			websocket.FormatCloseMessage(
				websocket.CloseUnsupportedData,
				"Invalid token",
			),
		)
		conn.Close()
		return
	}
	communicate.RegisterPlayerConnection(user.Email, conn)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
	communicate.UnregisterPlayerConnection(user.Email, conn)
}
