package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardFriendOfChildren struct{}

func (c CardFriendOfChildren) Name() string {
	return "Friend of Children"
}

func (c CardFriendOfChildren) Description() string {
	return fmt.Sprintf(
		"When you gain any %s {{b}}gain 1 extra %s.{{_b}}",
		EnergySymbol, EnergySymbol,
	)
}

func (c CardFriendOfChildren) Cost() int {
	return 3
}

func (c CardFriendOfChildren) Kind() int {
	return CardKindKeep
}

func (c CardFriendOfChildren) ModifyEnergy(game *Game, player, amount int) int {
	if amount > 0 {
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gains %s for gaining energy ({{b}}%s{{_b}})",
			game.RenderName(player),
			RenderEnergyChange(1),
			c.Name(),
		)))
		amount += 1
	}
	return amount
}
