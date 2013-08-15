package main

import (
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"strings"
)

type PokeCommand struct {
	gameId interface{}
}

func (pc PokeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("poke", 0, input)
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
	if pc.gameId != nil {
		CommunicateGameTo(pc.gameId, g, whoseTurn, fmt.Sprintf(
			"%s wants to remind you it's your turn!", player), false)
	}
	return "You poked " + strings.Join(whoseTurn, ", "), nil
}

func (pc PokeCommand) Usage(player string, context interface{}) string {
	return "{{b}}poke{{_b}} to remind the other players to play"
}
