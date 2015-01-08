package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CardCompleteDestruction struct{}

func (c CardCompleteDestruction) Name() string {
	return "Complete Destruction"
}

func (c CardCompleteDestruction) Description() string {
	return fmt.Sprintf(
		"If you roll %s%s%s%s%s%s {{b}}gain %s{{_b}} in addition to the regular results.",
		RenderDie(Die1),
		RenderDie(Die2),
		RenderDie(Die3),
		RenderDie(DieHeal),
		RenderDie(DieAttack),
		RenderDie(DieEnergy),
		RenderVP(9),
	)
}

func (c CardCompleteDestruction) Cost() int {
	return 3
}

func (c CardCompleteDestruction) Kind() int {
	return CardKindKeep
}

func (c CardCompleteDestruction) HandlePreResolveDice(
	game *Game,
	player int,
	dice []int,
) []int {
	rolledDice := map[int]bool{}
	for _, d := range dice {
		rolledDice[d] = true
	}
	if len(rolledDice) == 6 {
		game.Boards[player].ModifyVP(9)
		game.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained %s ({{b}}%s{{_b}})",
			game.RenderName(player),
			RenderVP(9),
			c.Name(),
		)))
	}
	return dice
}
