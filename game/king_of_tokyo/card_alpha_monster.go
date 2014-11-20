package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardAlphaMonster struct{}

func (c CardAlphaMonster) Name() string {
	return "Alpha Monster"
}

func (c CardAlphaMonster) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} when you attack.",
		RenderVP(1),
	)
}

func (c CardAlphaMonster) Cost() int {
	return 5
}

func (c CardAlphaMonster) Kind() int {
	return CardKindKeep
}

func (c CardAlphaMonster) PostAttack(game *Game, player, attack int) {
	game.Boards[player].VP += 1
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s for attacking ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(1),
		c.Name(),
	)))
}
