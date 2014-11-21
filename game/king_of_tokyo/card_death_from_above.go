package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardDeathFromAbove struct{}

func (c CardDeathFromAbove) Name() string {
	return "Death From Above"
}

func (c CardDeathFromAbove) Description() string {
	return fmt.Sprintf(
		"%s {{b}}and take control of Tokyo{{_b}} if you don't already control it",
		RenderVPChange(2),
	)
}

func (c CardDeathFromAbove) Cost() int {
	return 5
}

func (c CardDeathFromAbove) Kind() int {
	return CardKindDiscard
}

func (c CardDeathFromAbove) PostCardBuy(game *Game, player int, card CardBase, cost int) {
	game.Boards[player].VP += 2
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(2),
		c.Name(),
	)))
	game.TakeControl(player)
}
