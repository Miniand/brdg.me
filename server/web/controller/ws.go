package controller

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
	"strings"
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
		} else {
			fmt.Printf("Handling: %s\n", string(in))
			HandleMessage(ws, email, string(in))
		}
	}
	fmt.Println("the end")
}

func HandleMessage(ws *websocket.Conn, email, message string) error {
	parts := strings.SplitN(message, ";", 2)
	if len(parts) != 2 {
		return nil
	}
	gm, err := model.LoadGame(parts[0])
	if err != nil {
		SendError(ws, "Could not find game")
		return nil
	}
	g, err := gm.ToGame()
	if err != nil {
		SendError(ws, "Could not load game")
		return nil
	}
	out, err := command.CallInCommands(email, g, parts[1], g.Commands())
	if out != "" {
		SendOutput(ws, out)
	}
	if err != nil {
		SendError(ws, err.Error())
		return nil
	}
	gm, err = model.UpdateGame(gm.Id, g)
	if err != nil {
		SendError(ws, "Could not save game")
		return nil
	}
	rawGameOutput, err := g.RenderForPlayer(email)
	if err != nil {
		SendError(ws, "Could not render game")
		return nil
	}
	gameOutput, err := render.RenderHtml(rawGameOutput)
	if err != nil {
		SendError(ws, "Could not render game")
		return nil
	}
	SendCommand(ws, "gameOutput", gameOutput)
	logOutput, err := render.RenderHtml(
		log.RenderMessages(g.GameLog().MessagesFor(email)))
	if err != nil {
		SendError(ws, "Could not render log")
		return nil
	}
	SendCommand(ws, "log", logOutput)
	commandOutput, err := render.RenderHtml(render.CommandUsages(
		command.CommandUsages(email, g,
			command.AvailableCommands(email, g, g.Commands()))))
	if err != nil {
		SendError(ws, "Could not render commands")
		return nil
	}
	SendCommand(ws, "commands", commandOutput)
	return nil
}

func SendCommand(ws *websocket.Conn, cmd, arg string) {
	websocket.Message.Send(ws, fmt.Sprintf("%s;%s", cmd, arg))
}

func SendOutput(ws *websocket.Conn, out string) {
	SendCommand(ws, "out", out)
}

func SendError(ws *websocket.Conn, err string) {
	SendCommand(ws, "err", err)
}
