package game

import (
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/web/view"
	"html/template"
	"io"
)

type ShowScope struct {
	GameModel *model.GameModel
	Game      game.Playable
	Output    template.HTML
}

var showTmpl = `{{template "header" "Game"}}
<div class="pure-g-r">
	<div class="pure-u-1-2 game-output-container">
		<div class="game-output">{{.Output}}</div>
	</div>
	<div class="pure-u-1-2 game-log"></div>
</div>
{{template "footer"}}`

func Show(wr io.Writer, scope ShowScope) {
	var err error
	scope.Game, err = scope.GameModel.ToGame()
	if err != nil {
		panic(err.Error())
	}
	rawOutput, err := scope.Game.RenderForPlayer("Mick")
	if err != nil {
		panic(err.Error())
	}
	output, err := render.RenderHtml(rawOutput)
	if err != nil {
		panic(err.Error())
	}
	scope.Output = template.HTML(output)
	view.Parse("gameShow", showTmpl).Execute(wr, scope)
}
