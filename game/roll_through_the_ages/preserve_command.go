package roll_through_the_ages

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type PreserveCommand struct{}

func (c PreserveCommand) Name() string { return "preserve" }

func (c PreserveCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	return "", g.Preserve(pNum)
}

func (c PreserveCommand) Usage(player string, context interface{}) string {
	return "{{b}}preserve{{_b}} to use 1 pottery to double your food"
}

func (g *Game) CanPreserve(player int) bool {
	b := g.Boards[player]
	return g.CurrentPlayer == player && g.Phase == PhasePreserve &&
		b.Developments[DevelopmentPreservation] && b.Goods[GoodPottery] > 0 &&
		b.Food > 0
}

func (g *Game) Preserve(player int) error {
	if !g.CanPreserve(player) {
		return errors.New("you can't preserve at the moment")
	}

	g.Boards[player].Food *= 2
	g.Boards[player].Goods[GoodPottery] -= 1

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s used {{b}}preservation{{_b}} to double their food to {{b}}%d{{_b}} for {{b}}1 pottery{{_b}}`,
		g.RenderName(player),
		g.Boards[player].Food,
	)))
	g.NextPhase()
	return nil
}
