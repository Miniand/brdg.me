package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type BuildCommand struct{}

func (c BuildCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("build", 1, -1, input)
}

func (c BuildCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanBuild(p)
}

func (c BuildCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	b, err := ParseBuildable(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Build(p, b)
}

func (c BuildCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(
		`{{b}}build ##{{_b}} to build a trade ship, colony ship, cannon or booster.  Eg. {{b}}build colony{{_b}}
{{b}}Item{{_b}}         {{b}}Price{{_b}}
trade ship   %s
colony ship  %s
cannon       %s
booster      %s`,
		TradeShipTransaction.LoseString(),
		ColonyShipTransaction.LoseString(),
		g.PlayerBoards[p].CannonTransaction().LoseString(),
		g.PlayerBoards[p].BoosterTransaction().LoseString(),
	)
}

func (g *Game) CanBuild(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseTradeAndBuild &&
		g.PlayerBoards[player].CanBuild()
}

func (g *Game) Build(player, resource int) error {
	if !g.CanBuild(player) {
		return errors.New("you cannot build at the moment")
	}
	var t Transaction
	switch resource {
	case ResourceTradeShip:
		t = TradeShipTransaction
	case ResourceColonyShip:
		t = ColonyShipTransaction
	case ResourceCannon:
		t = g.PlayerBoards[player].CannonTransaction()
	case ResourceBooster:
		t = g.PlayerBoards[player].BoosterTransaction()
	default:
		return errors.New("invalid resource")
	}
	if !g.PlayerBoards[player].CanAfford(t) {
		return t.CannotAffordError()
	}
	if !g.PlayerBoards[player].CanFit(t) {
		return t.CannotFitError()
	}
	g.PlayerBoards[player].Transact(t)
	g.LogTransaction(player, t)
	return nil
}
