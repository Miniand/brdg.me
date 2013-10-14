package game

import (
	"github.com/Miniand/brdg.me/server/web/view"
	"io"
)

var indexTmpl = `{{template "header" "Games"}}
Games
{{template "footer"}}`

func Index(wr io.Writer) {
	view.Parse("gameIndex", indexTmpl).Execute(wr, nil)
}
