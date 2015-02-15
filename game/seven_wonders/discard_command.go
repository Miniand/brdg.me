package seven_wonders

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
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
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which numbered card to discard")
	}
	cardNum, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("that is not a valid card number")
	}
	return "", g.DiscardCard(pNum, cardNum-1)
}

func (c DiscardCommand) Usage(player string, context interface{}) string {
	return fmt.Sprintf(
		"{{b}}discard #{{_b}} to discard a card for %s, eg. {{b}}discard 2{{_b}}",
		RenderMoney(3),
	)
}

func (g *Game) CanDiscard(player int) bool {
	return g.CanAction(player)
}

func (g *Game) DiscardCard(player, cardNum int) error {
	if !g.CanDiscard(player) {
		return errors.New("cannot discard at the moment")
	}
	if cardNum < 0 || cardNum >= len(g.Hands[player]) {
		return errors.New("not a valid card number")
	}
	g.Actions[player] = &DiscardAction{
		Card: cardNum,
	}
	g.CheckHandComplete()
	return nil
}
