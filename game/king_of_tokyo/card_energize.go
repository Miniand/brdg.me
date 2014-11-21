package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardEnergize struct{}

func (c CardEnergize) Name() string {
	return "Energize"
}

func (c CardEnergize) Description() string {
	return RenderEnergyChange(9)
}

func (c CardEnergize) Cost() int {
	return 8
}

func (c CardEnergize) Kind() int {
	return CardKindDiscard
}

func (c CardEnergize) PostCardBuy(game *Game, player int, card CardBase, cost int) {
	game.Boards[player].ModifyEnergy(9)
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderEnergy(1),
		c.Name(),
	)))
}
