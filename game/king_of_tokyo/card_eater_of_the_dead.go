package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardEaterOfTheDead struct{}

func (c CardEaterOfTheDead) Name() string {
	return "Eater of the Dead"
}

func (c CardEaterOfTheDead) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} every time a monster's %s goes to 0.",
		RenderVP(3),
		HealthSymbol,
	)
}

func (c CardEaterOfTheDead) Cost() int {
	return 4
}

func (c CardEaterOfTheDead) Kind() int {
	return CardKindKeep
}

func (c CardEaterOfTheDead) HealthZero(game *Game, player, zeroPlayer int) {
	game.Boards[player].VP += 3
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderVP(3),
		c.Name(),
	)))
}
