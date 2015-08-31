package red7

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type PlayCommand struct{}

func (c PlayCommand) Name() string { return "play" }

func (c PlayCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("you must specify one card")
	}
	card, ok := ParseCard(args[0])
	if !ok {
		return "", errors.New("the card must be a letter followed by a number, eg. r6")
	}

	return "", g.Play(pNum, card)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a card to your palette, eg. {{b}}play b4{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return g.CurrentPlayer == player && !g.HasPlayed
}

func (g *Game) Play(player, card int) error {
	if !g.CanPlay(player) {
		return errors.New("you can't play at the moment")
	}
	index, ok := helper.IntFind(card, g.Hands[player])
	if !ok {
		return errors.New("you don't have that card")
	}
	g.Palettes[player] = append(g.Palettes[player], card)
	g.Hands[player] = append(g.Hands[player][:index], g.Hands[player][index+1:]...)
	g.HasPlayed = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s",
		g.PlayerName(player),
		RenderCard(card),
	)))
	return nil
}
