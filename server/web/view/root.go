package view

import (
	"github.com/Miniand/brdg.me/game"
	"io"
)

type rootScope struct {
	Games map[string]game.Playable
}

var rootTmpl = `{{template "header" "Play board games online by web and email"}}
<div class="title">
	<h1>{{template "title"}}

Play board games online by web and email</h1>
</div>
<div class="pure-g-r game-list">
	{{range .Games}}<div class="pure-u-1-3 game">
		<h2><a href="/game/new/{{.Identifier}}">{{.Name}}</a></h2>
		<p>Hotels and shit</p>
	</div>{{end}}
</div>
{{template "footer"}}`

func Root(wr io.Writer) {
	Parse("root", rootTmpl).Execute(wr, rootScope{
		Games: game.RawCollection(),
	})
}
