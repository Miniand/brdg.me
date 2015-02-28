package communicate

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"

	"github.com/gorilla/websocket"
)

const (
	WsTypeGameUpdate = "gameUpdate"
)

type WsMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WsGameUpdateMsg struct {
	Text     string `json:"text,omitempty"`
	TextHtml string `json:"textHtml,omitempty"`
	MsgType  string `json:"msgType,omitempty"`
	GameId   string `json:"gameId"`
	GameName string `json:"gameName"`
	YourTurn bool   `json:"yourTurn"`
}

var wsConnections = map[string][]*websocket.Conn{}

var errNoConnections = errors.New(
	"that player does not have any connections open")

func NewWsGameUpdateMsg(
	player, text, textHtml, msgType string,
	g game.Playable,
	gm *model.GameModel,
) WsGameUpdateMsg {
	isFinished := gm.IsFinished
	yourTurn := false
	if !isFinished {
		for _, p := range gm.WhoseTurn {
			if p == player {
				yourTurn = true
				break
			}
		}
	}
	return WsGameUpdateMsg{
		Text:     text,
		TextHtml: textHtml,
		MsgType:  msgType,
		GameId:   gm.Id,
		GameName: g.Name(),
		YourTurn: yourTurn,
	}
}

func wsSendGameMulti(
	players []string,
	text, msgType string,
	g game.Playable,
	gm *model.GameModel,
) (
	failed map[string]error) {
	failed = map[string]error{}
	for _, p := range players {
		if err := wsSendGame(p, text, msgType, g, gm); err != nil {
			failed[p] = err
		}
	}
	return
}

func wsSendGame(
	player, text, msgType string,
	g game.Playable,
	gm *model.GameModel,
) (err error) {
	textHtml, err := render.RenderHtml(text)
	if err != nil {
		return fmt.Errorf("unable to render text to HTML: %v", err)
	}
	return wsSend(player, WsTypeGameUpdate, NewWsGameUpdateMsg(
		player,
		text,
		textHtml,
		msgType,
		g,
		gm,
	))
}

func wsSend(player string, msgType string, data interface{}) (err error) {
	sent := false
	conns := wsConnections[player]
	if conns == nil || len(conns) == 0 {
		return errNoConnections
	}
	for _, conn := range conns {
		if err = conn.WriteJSON(WsMsg{
			Type: msgType,
			Data: data,
		}); err == nil {
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
