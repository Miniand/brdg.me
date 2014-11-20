package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardAlienMetabolism struct{}

func (c CardAlienMetabolism) Name() string {
	return "Alien Metabolism"
}

func (c CardAlienMetabolism) Description() string {
	return fmt.Sprintf(
		`{{b}}Buying cards costs you 1 less %s.{{_b}}`,
		EnergySymbol,
	)
}

func (c CardAlienMetabolism) Cost() int {
	return 3
}

func (c CardAlienMetabolism) Kind() int {
	return CardKindKeep
}

func (c CardAlienMetabolism) ModifyCardCost(game *Game, player, cost int) int {
	cost = cost - 1
	if cost < 0 {
		cost = 0
	}
	return cost
}

func (c CardAlienMetabolism) PostCardBuy(game *Game, card CardBase, cost int) {
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Card cost 1 less ({{b}}%s{{_b}})",
		c.Name(),
	)))
}
