package scommand

import (
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

func CommandsForGame(gm *model.GameModel, g game.Playable) []comm.Command {
	c := []comm.Command{}
	if !IsConcedeVoting(gm) {
		c = g.Commands()
	}
	c = append(c, Commands(gm)...)
	return c
}

func Commands(gm *model.GameModel) []comm.Command {
	return []comm.Command{
		ConcedeCommand{
			gameModel: gm,
		},
		ConcedeYesCommand{
			gameModel: gm,
		},
		ConcedeNoCommand{
			gameModel: gm,
		},
		PokeCommand{
			gameModel: gm,
		},
		SayCommand{
			gameModel: gm,
		},
		NewCommand{},
		RestartCommand{
			gameModel: gm,
		},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}
}
