package game

import (
	"github.com/Miniand/brdg.me/server/web/view"
	"io"
)

var showTmpl = `{{template "header" "Game"}}
Game {{.Game.Id}}
{{template "footer"}}`

func Show(wr io.Writer) {
	view.Parse("gameShow", showTmpl).Execute(wr, nil)
}
