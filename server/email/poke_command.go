package main

import (
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
)

type PokeCommand struct{}

func (pc PokeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("poke", 0, input)
}

func (pc PokeCommand) CanCall(player string, context interface{}) bool {
	if currentGameId == nil {
		return false
	}
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

func (pc PokeCommand) Call(player string, context interface{}, args []string) error {
	if currentGameId == nil {
		return errors.New("There is no relevant game ID to poke for")
	}
	g, ok := context.(game.Playable)
	if !ok || g.IsFinished() {
		return errors.New("The game is already finished")
	}
	CommunicateGameTo(currentGameId, g, g.WhoseTurn(),
		fmt.Sprintf("%s wants to remind you it's your turn!", player), false)
	return nil
}

func (pc PokeCommand) Usage(player string, context interface{}) string {
	return "{{b}}poke{{_b}} to remind the other players to play"
}
