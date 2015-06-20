package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardHerbivore struct{}

func (c CardHerbivore) Name() string {
	return "Herbivore"
}

func (c CardHerbivore) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} on your turn if you don't damage anyone.",
		RenderVP(1),
	)
}

func (c CardHerbivore) Cost() int {
	return 5
}

func (c CardHerbivore) Kind() int {
	return CardKindKeep
}

func (c *CardHerbivore) HandleStartTurn(game *Game, player int) {
	game.Boards[player].HasDealtDamage = false
}

func (c *CardHerbivore) HandleDamageDealt(
	game *Game,
	player, target, damage int,
) {
	game.Boards[player].HasDealtDamage = true
}

func (c CardHerbivore) HandleEndTurn(game *Game, player int) {
	if !game.Boards[player].HasDealtDamage {
		game.Boards[player].ModifyVP(1)
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained %s ({{b}}%s{{_b}})",
			game.RenderName(player),
			RenderVP(1),
			c.Name(),
		)))
	}
}
