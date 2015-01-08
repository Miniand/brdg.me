package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardFreezeTime struct{}

func (c CardFreezeTime) Name() string {
	return "Freeze Time"
}

func (c CardFreezeTime) Description() string {
	oneDie := RenderDie(Die1)
	return fmt.Sprintf(
		"On a turn where you score %s%s%s, {{b}}you can take another turn{{_b}} with one less die.",
		oneDie, oneDie, oneDie,
	)
}

func (c CardFreezeTime) Cost() int {
	return 5
}

func (c CardFreezeTime) Kind() int {
	return CardKindKeep
}

func (c CardFreezeTime) HandlePreResolveDice(
	game *Game,
	player int,
	dice []int,
) []int {
	count := 0
	for _, d := range dice {
		if d == Die1 {
			count += 1
		}
	}
	if count >= 3 {
		extraTurnDice := len(game.CurrentRoll) - 1
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s will get an extra turn with %d dice ({{b}}%s{{_b}})",
			game.RenderName(player),
			extraTurnDice,
			c.Name(),
		)))
		game.ExtraTurns = append(game.ExtraTurns, extraTurnDice)
	}
	return dice
}
