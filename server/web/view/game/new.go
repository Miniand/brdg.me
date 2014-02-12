package game

import (
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/web/view"
	"io"
)

type NewScope struct {
	Game    game.Playable
	Players []string
}

var newTmpl = `{{template "header" "New game"}}
<h1>New game of {{.Game.Name}}</h1>
<h2>Players</h2>
<form method="POST" action="/game">
	<input type="hidden" name="identifier" value="{{.Game.Identifier}}" />
	{{loggedInUser}}
	{{range .Players}}
		<input name="players[]" value="{{.}}" />
	{{end}}
	<input name="players[]" />
	<input type="button" onclick="$(this).before($('<input/>').attr('name', 'players[]'))" value="Add" />
	<input type="submit" value="Start game!" />
</form>
{{template "footer"}}`

func New(wr io.Writer, scope NewScope) {
	view.Parse("gameNew", newTmpl).Execute(wr, scope)
}
