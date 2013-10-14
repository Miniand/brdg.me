package game

import (
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/web/view"
	"io"
)

type NewScope struct {
	Game game.Playable
}

var newTmpl = `{{template "header" "New game"}}
Game {{.Game.Name}}
{{template "footer"}}`

func New(wr io.Writer, scope NewScope) {
	view.Parse("gameNew", newTmpl).Execute(wr, scope)
}
