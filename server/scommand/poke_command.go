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

const (
	MsgTypePoke = "poke"
)

type PokeCommand struct {
	gameModel *model.GameModel
}

func (pc PokeCommand) Parse(input string) []string {
	return comm.ParseNamedCommandNArgs("poke", 0, input)
}

func (pc PokeCommand) CanCall(player string, context interface{}) bool {
	if pc.gameModel.IsFinished {
		return false
	}
	waitingOn := pc.gameModel.WhoseTurn
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
	if pc.gameModel.IsFinished {
		return "", errors.New("The game is already finished")
	}
	whoseTurn := pc.gameModel.WhoseTurn
	if pc.gameModel != nil && pc.gameModel.Id != "" {
		communicate.Game(
			g,
			pc.gameModel,
			whoseTurn,
			CommandsForGame(pc.gameModel, g),
			fmt.Sprintf(
				"%s wants to remind you it's your turn!",
				render.PlayerNameInPlayers(player, pc.gameModel.PlayerList)),
			MsgTypePoke,
			false,
		)
	}
	return "You poked the other players to take their turn", nil
}

func (pc PokeCommand) Usage(player string, context interface{}) string {
	return "{{b}}poke{{_b}} to remind the other players to play"
}
