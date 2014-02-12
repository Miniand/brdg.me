package game

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/web/view"
	"html/template"
	"io"
)

type ShowScope struct {
	GameModel         *model.GameModel
	Game              game.Playable
	Output            template.HTML
	Log               template.HTML
	AvailableCommands template.HTML
}

var showTmpl = `{{template "header" "Game"}}
<div class="game-show">
	<div class="pure-g-r game-show">
		<div class="pure-u-1-2 game-output-container">
			<h2 class="game-output-heading">Game</h2>
			<div class="game-output">{{.Output}}</div>
		</div>
		<div class="pure-u-1-2 game-log-container">
			<h2 class="game-log-heading">Log</h2>
			<div class="game-log">{{.Log}}</div>
		</div>
	</div>
	<div class="game-input-container">
		<h2 class="game-input-heading">You can:</h2>
		<div class="game-input-available-commands">{{.AvailableCommands}}</div>
		<form class="game-input" onsubmit="var i = $('.game-input-command');ws.send('{{.GameModel.Id}};'+i.val());i.val('');return false;">
			<input type="text" class="game-input-command" />
			<input type="submit" class="game-input-submit" value="Send command" />
		</form>
	</div>
</div>
{{template "footer"}}
<script>
function logScrollToBottom() {
	var lc = document.querySelector(".game-log");
	lc.scrollTop = lc.scrollHeight;	
}
logScrollToBottom();
</script>`

func Show(wr io.Writer, scope ShowScope) {
	var err error
	// Game
	scope.Game, err = scope.GameModel.ToGame()
	if err != nil {
		panic(err.Error())
	}
	// Output
	rawOutput, err := scope.Game.RenderForPlayer(view.LoggedInUser)
	if err != nil {
		panic(err.Error())
	}
	output, err := render.RenderHtml(rawOutput)
	if err != nil {
		panic(err.Error())
	}
	scope.Output = template.HTML(output)
	// Log
	rawLog := log.RenderMessages(scope.Game.GameLog().MessagesFor(
		view.LoggedInUser))
	logOutput, err := render.RenderHtml(rawLog)
	if err != nil {
		panic(err.Error())
	}
	scope.Log = template.HTML(logOutput)
	// Available commands
	rawCommands := render.CommandUsages(command.CommandUsages(
		view.LoggedInUser, scope.Game, command.AvailableCommands(
			view.LoggedInUser, scope.Game, scope.Game.Commands())))
	commandOutput, err := render.RenderHtml(rawCommands)
	if err != nil {
		panic(err.Error())
	}
	scope.AvailableCommands = template.HTML(commandOutput)
	view.Parse("gameShow", showTmpl).Execute(wr, scope)
}
