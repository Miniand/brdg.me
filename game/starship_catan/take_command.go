package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type TakeCommand struct{}

func (c TakeCommand) Name() string { return "take" }

func (c TakeCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)

	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", errors.New("could not parse player")
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("you must pass a resource argument")
	}

	r, err := helper.MatchStringInStringMap(args[0], ResourceNameMap(Goods))
	if err != nil {
		return "", err
	}

	return "", g.Take(p, r)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return fmt.Sprintf(
		"{{b}}take ##{{_b}} to take a good from your opponent for %s, eg. {{b}}take carbon{{_b}}",
		RenderMoney(2),
	)
}

func TakeTransaction(resource int) Transaction {
	return Transaction{
		resource:      1,
		ResourceAstro: -2,
	}
}

func (g *Game) CanTake(player int) bool {
	return g.CanTakeResource(player, ResourceAny)
}

func (g *Game) CanTakeResource(player, resource int) bool {
	if !(g.CurrentPlayer == player && g.Phase == PhaseTradeAndBuild &&
		g.RemainingPlayerTrades() > 0) {
		return false
	}
	if resource == ResourceAny {
		return true
	}
	opponent := (player + 1) % 2
	t := TakeTransaction(resource)
	return g.PlayerBoards[player].CanFit(t) &&
		g.PlayerBoards[opponent].CanFit(t.Inverse())
}

func (g *Game) Take(player, resource int) error {
	if !g.CanTakeResource(player, resource) {
		return errors.New("can't take that resource at the moment")
	}
	if !ContainsInt(resource, Goods) {
		return errors.New("you can only take goods")
	}
	opponent := (player + 1) % 2
	t := TakeTransaction(resource)
	g.PlayerBoards[player].Transact(t)
	g.LogTransaction(player, t)
	g.PlayerBoards[opponent].Transact(t.Inverse())
	g.LogTransaction(opponent, t.Inverse())
	g.PlayerTradeAmount += 1
	return nil
}
