package communicate

import (
	"errors"
	"reflect"

	"github.com/Miniand/brdg.me/game"

	"github.com/gorilla/websocket"
)

type WsMsg struct {
	Text     string `json:"text,omitempty"`
	GameId   string `json:"gameId"`
	GameName string `json:"gameName"`
	YourTurn bool   `json:"yourTurn"`
}

var wsConnections = map[string][]*websocket.Conn{}

var errNoConnections = errors.New(
	"that player does not have any connections open")

func NewWsMsg(player, gameId, text string, g game.Playable) WsMsg {
	isFinished := g.IsFinished()
	yourTurn := false
	if !isFinished {
		for _, p := range g.WhoseTurn() {
			if p == player {
				yourTurn = true
				break
			}
		}
	}
	return WsMsg{
		Text:     text,
		GameId:   gameId,
		GameName: g.Name(),
		YourTurn: yourTurn,
	}
}

func wsSendGameMulti(players []string, gameId, text string, g game.Playable) (
	failed map[string]error) {
	failed = map[string]error{}
	for _, p := range players {
		if err := wsSendGame(p, gameId, text, g); err != nil {
			failed[p] = err
		}
	}
	return
}

func wsSendGame(player, gameId, text string, g game.Playable) (err error) {
	sent := false
	conns := wsConnections[player]
	if conns == nil || len(conns) == 0 {
		return errNoConnections
	}
	for _, conn := range conns {
		if err = conn.WriteJSON(NewWsMsg(
			player,
			gameId,
			text,
			g,
		)); err == nil {
			sent = true
		}
	}
	if !sent {
		if err == nil {
			err = errNoConnections
		}
	} else {
		err = nil
	}
	return
}

func RegisterPlayerConnection(player string, conn *websocket.Conn) {
	if wsConnections[player] == nil {
		wsConnections[player] = []*websocket.Conn{}
	}
	wsConnections[player] = append(wsConnections[player], conn)
}

func UnregisterPlayerConnection(player string, conn *websocket.Conn) {
	cval := reflect.ValueOf(conn).Pointer()
	if wsConnections[player] == nil {
		return
	}
	for i, c := range wsConnections[player] {
		if reflect.ValueOf(c).Pointer() == cval {
			wsConnections[player] = append(
				wsConnections[player][:i], wsConnections[player][i+1:]...)
			return
		}
	}
}