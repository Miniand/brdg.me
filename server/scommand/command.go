package scommand

import (
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

func CommandsForGame(
	player string,
	gm *model.GameModel,
	g game.Playable,
) []comm.Command {
	c := []comm.Command{}
	if !gm.IsFinished && !gm.IsConcedeVoting() {
		c = g.Commands(player)
	}
	c = append(c, Commands(player, gm)...)
	return c
}

func Commands(player string, gm *model.GameModel) []comm.Command {
	c := []comm.Command{}
	if gm != nil {
		if CanConcedeVote(player, gm) {
			c = append(
				c,
				ConcedeVoteCommand{
					gameModel: gm,
				},
			)
		}
		if CanSay(player, gm) {
			c = append(c, SayCommand{
				gameModel: gm,
			})
		}
		if CanInitiateConcedeVote(player, gm) {
			c = append(c, ConcedeCommand{
				gameModel: gm,
			})
		}
		if CanRestart(player, gm) {
			c = append(c, RestartCommand{
				gameModel: gm,
			})
		}
		if CanDump(player, gm) {
			c = append(c, DumpCommand{
				gameModel: gm,
			})
		}
	}
	if CanNew(player, gm) {
		c = append(c, NewCommand{})
	}
	if CanList(player, gm) {
		c = append(c, ListCommand{})
	}
	if CanUnsubscribe(player, gm) {
		c = append(c, UnsubscribeCommand{})
	}
	if CanSubscribe(player, gm) {
		c = append(c, SubscribeCommand{})
	}
	return c
}
