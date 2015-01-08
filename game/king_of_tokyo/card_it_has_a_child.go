package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardItHasAChild struct{}

func (c CardItHasAChild) Name() string {
	return "It Has a Child!"
}

func (c CardItHasAChild) Description() string {
	return fmt.Sprintf(
		"If you are eliminated discard all your cards and lose all your %s. {{b}}Heal to %s and start again.{{_b}}",
		VPSymbol,
		RenderHealth(10),
	)
}

func (c CardItHasAChild) Cost() int {
	return 7
}

func (c CardItHasAChild) Kind() int {
	return CardKindKeep
}

func (c CardItHasAChild) HandleHealthZero(game *Game, player, zeroPlayer int) {
	if player == zeroPlayer {
		game.Discard = append(game.Discard, game.Boards[player].Cards...)
		game.Boards[player].Cards = []CardBase{}
		game.Boards[player].Health = 10
		game.Boards[player].VP = 0
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s has been revived and has lost all their cards and VP ({{b}}%s{{_b}})",
			game.RenderName(player),
			c.Name(),
		)))
	}
}
