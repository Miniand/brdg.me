package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardEvacuationOrders struct{}

func (c CardEvacuationOrders) Name() string {
	return "Evacuation Orders"
}

func (c CardEvacuationOrders) Description() string {
	return fmt.Sprintf(
		"{{b}}All other monsters lose %s{{_b}}.",
		RenderVP(5),
	)
}

func (c CardEvacuationOrders) Cost() int {
	return 7
}

func (c CardEvacuationOrders) Kind() int {
	return CardKindDiscard
}

func (c CardEvacuationOrders) PostCardBuy(game *Game, player int, card CardBase, cost int) {
	for p, _ := range game.Players {
		if p == player {
			continue
		}
		game.Boards[p].ModifyVP(-5)
	}
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s made all other monsters lose %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(5),
		c.Name(),
	)))
}
