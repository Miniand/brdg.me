package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type BuildCommand struct{}

func (c BuildCommand) Name() string { return "build" }

func (c BuildCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify something to build")
	}
	b, err := ParseBuildable(args[0])
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
