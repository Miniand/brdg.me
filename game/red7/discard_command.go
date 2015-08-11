package red7

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type DiscardCommand struct{}

func (c DiscardCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("discard", 1, input)
}

func (c DiscardCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanDiscard(pNum)
}

func (c DiscardCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) != 1 {
		return "", errors.New("you must specify one card")
	}
	card, ok := ParseCard(a[0])
	if !ok {
		return "", errors.New("the card must be a letter followed by a number, eg. r6")
	}

	return "", g.Discard(pNum, card)
}

func (c DiscardCommand) Usage(player string, context interface{}) string {
	return "{{b}}discard ##{{_b}} to discard a card and set the new rule, eg. {{b}}discard b4{{_b}}"
}

func (g *Game) CanDiscard(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Discard(player, card int) error {
	if !g.CanDiscard(player) {
		return errors.New("you can't discard at the moment")
	}
	index, ok := helper.IntFind(card, g.Hands[player])
	if !ok {
		return errors.New("you don't have that card")
	}
	if leader, _ := g.LeaderWithSuit(card & SuitMask); leader != player {
		return errors.New("you wouldn't be the leader after discarding that card")
	}
	g.DiscardPile = append(g.DiscardPile, card)
	g.Hands[player] = append(g.Hands[player][:index], g.Hands[player][index+1:]...)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded %s, the new rule is {{b}}%s{{_b}}",
		g.PlayerName(player),
		RenderCard(card),
		SuitRulesStrs[card&SuitMask],
	)))
	if RankVal[card&RankMask] > len(g.Palettes[player]) {
		// Draw a card for discarding a card numbered higher than palette size.
		g.Draw(player, 1)
	}
	g.EndTurn()
	return nil
}
