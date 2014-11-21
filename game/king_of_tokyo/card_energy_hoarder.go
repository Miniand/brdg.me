package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardEnergyHoarder struct{}

func (c CardEnergyHoarder) Name() string {
	return "Energy Hoarder"
}

func (c CardEnergyHoarder) Description() string {
	return fmt.Sprintf(
		"{{b}}You gain %s{{_b}} for every %s you have at the end of your turn.",
		RenderVP(1),
		RenderEnergy(6),
	)
}

func (c CardEnergyHoarder) Cost() int {
	return 3
}

func (c CardEnergyHoarder) Kind() int {
	return CardKindKeep
}

func (c CardEnergyHoarder) EndTurn(game *Game, player int) {
	if vp := game.Boards[player].Energy / 6; vp > 0 {
		game.Boards[player].ModifyVP(vp)
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained %s ({{b}}%s{{_b}})",
			game.RenderName(player),
			RenderVP(vp),
			c.Name(),
		)))
	}
}
