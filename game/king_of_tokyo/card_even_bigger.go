package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardEvenBigger struct{}

func (c CardEvenBigger) Name() string {
	return "Even Bigger"
}

func (c CardEvenBigger) Description() string {
	return fmt.Sprintf(
		"{{b}}Your maximum %s is increased by 2.  Gain %s{{_b}} when you get this card.",
		HealthSymbol,
		RenderHealth(2),
	)
}

func (c CardEvenBigger) Cost() int {
	return 4
}

func (c CardEvenBigger) Kind() int {
	return CardKindKeep
}

func (c CardEvenBigger) ModifyMaxHealth(health int) int {
	return health + 2
}

func (c CardEvenBigger) HandlePostCardBuy(
	game *Game,
	player int,
	card CardBase,
	cost int,
) {
	game.Boards[player].ModifyHealth(2)
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s gained %s ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderHealth(2),
		c.Name(),
	)))
}

func (c CardEvenBigger) HandleCardLost(game *Game, player int) {
	game.Boards[player].ModifyHealth(-2)
	game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s lost %s for losing a card ({{b}}%s{{_b}})",
		game.RenderName(player),
		RenderHealth(-2),
		c.Name(),
	)))
}
