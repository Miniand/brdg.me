package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardApartmentBuilding struct{}

func (c CardApartmentBuilding) Name() string {
	return "Apartment Building"
}

func (c CardApartmentBuilding) Description() string {
	return RenderVPChange(3)
}

func (c CardApartmentBuilding) Cost() int {
	return 5
}

func (c CardApartmentBuilding) Kind() int {
	return CardKindDiscard
}

func (c CardApartmentBuilding) PostCardBuy(game *Game, player int, card CardBase, cost int) {
	game.Boards[player].VP += 3
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(3),
		c.Name(),
	)))
}
