package scommand

import (
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

func CommandsForGame(gm *model.GameModel, g game.Playable) []comm.Command {
	c := []comm.Command{}
	if !gm.IsFinished && !gm.IsConcedeVoting() {
		c = g.Commands()
	}
	c = append(c, Commands(gm)...)
	return c
}

func Commands(gm *model.GameModel) []comm.Command {
	c := []comm.Command{}
	if gm != nil {
		c = []comm.Command{
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
			ConcedeCommand{
				gameModel: gm,
			},
			RestartCommand{
				gameModel: gm,
			},
		}
	}
	c = append(c, []comm.Command{
		NewCommand{},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}...)
	return c
}
