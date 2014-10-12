package scommand

import (
	"errors"
	"fmt"
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/communicate"
	"github.com/Miniand/brdg.me/server/model"
)

type PokeCommand struct {
	gameModel *model.GameModel
}

func (pc PokeCommand) Parse(input string) []string {
	return comm.ParseNamedCommandNArgs("poke", 0, input)
}

func (pc PokeCommand) CanCall(player string, context interface{}) bool {
	g, ok := context.(game.Playable)
	if !ok || g.IsFinished() {
		return false
	}
	waitingOn := g.WhoseTurn()
	for _, p := range waitingOn {
		if p == player {
			// No point poking yourself!
			return false
		}
	}
	return true
}

func (pc PokeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("No game was passed in")
	}
	if g.IsFinished() {
		return "", errors.New("The game is already finished")
	}
	whoseTurn := g.WhoseTurn()
	if pc.gameModel != nil && pc.gameModel.Id != "" {
		communicate.Game(pc.gameModel.Id, g, whoseTurn,
			append(g.Commands(), Commands(pc.gameModel)...),
			fmt.Sprintf(
				"%s wants to remind you it's your turn!",
				render.PlayerNameInPlayers(player, g.PlayerList())), false)
	}
	return "You poked the current turn players", nil
}

func (pc PokeCommand) Usage(player string, context interface{}) string {
	return "{{b}}poke{{_b}} to remind the other players to play"
}
