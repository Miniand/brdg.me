package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardBurrowing struct{}

func (c CardBurrowing) Name() string {
	return "Burrowing"
}

func (c CardBurrowing) Description() string {
	return "{{b}}Deal 1 extra damage when in Tokyo. Deal 1 damage when yielding Tokyo{{_b}} to the monster taking it."
}

func (c CardBurrowing) Cost() int {
	return 5
}

func (c CardBurrowing) Kind() int {
	return CardKindKeep
}

func (c CardBurrowing) ModifyAttack(
	game *Game,
	player, damage int,
	attacked []int,
) (int, []int) {
	if game.PlayerLocation(game.CurrentPlayer) != LocationOutside {
		damage += 1
	}
	return damage, attacked
}

func (c CardBurrowing) LeaveTokyo(game *Game, location, player, enteringPlayer int) {
	if enteringPlayer != TokyoEmpty {
		game.TakeDamage(enteringPlayer, 1)
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s took 1 damage for taking Tokyo from %s ({{b}}%s{{_b}})",
			game.RenderName(enteringPlayer),
			game.RenderName(player),
			c.Name(),
		)))
	}
}
