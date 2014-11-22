package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardGourmet struct{}

func (c CardGourmet) Name() string {
	return "Gourmet"
}

func (c CardGourmet) Description() string {
	oneDie := RenderDie(Die1)
	return fmt.Sprintf(
		"When scoring %s%s%s, {{b}}gain 2 extra %s.{{_b}}",
		oneDie, oneDie, oneDie,
		VPSymbol,
	)
}

func (c CardGourmet) Cost() int {
	return 4
}

func (c CardGourmet) Kind() int {
	return CardKindKeep
}

func (c CardGourmet) PreResolveDice(game *Game, player int, dice []int) []int {
	count := 0
	for _, d := range dice {
		if d == Die1 {
			count += 1
		}
	}
	if count >= 3 {
		game.Boards[player].ModifyVP(2)
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained %s ({{b}}%s{{_b}})",
			game.RenderName(player),
			RenderVP(2),
			c.Name(),
		)))
	}
	return dice
}
