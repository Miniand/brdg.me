package controller

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

func Ws(ws *websocket.Conn) {
	defer ws.Close()
	email := ""
	for {
		var in []byte
		if err := websocket.Message.Receive(ws, &in); err != nil {
			break
		}
		if email == "" {
			if err := secCookie.Decode("session", string(in), &email); err != nil {
				break
			}
			fmt.Println("Authenticated as " + email)
			if err := websocket.Message.Send(ws, "You have authenticated as "+
				email); err != nil {
				break
			}
		} else {
			fmt.Println(email + " says: " + string(in))
			if err := websocket.Message.Send(ws, "You said: "+
				string(in)); err != nil {
				break
			}
		}
	}
	fmt.Println("the end")
}
