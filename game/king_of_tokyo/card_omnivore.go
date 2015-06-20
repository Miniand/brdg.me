package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
)

type CardOmnivore struct{}

func (c CardOmnivore) Name() string {
	return "Omnivore"
}

func (c CardOmnivore) Description() string {
	return fmt.Sprintf(
		"{{b}}Once on each turn you can score %s%s%s for %s.{{_b}} You can also use these dice in other combinations.",
		RenderDie(Die1),
		RenderDie(Die2),
		RenderDie(Die3),
		RenderVP(2),
	)
}

func (c CardOmnivore) Cost() int {
	return 4
}

func (c CardOmnivore) Kind() int {
	return CardKindKeep
}

func (c CardOmnivore) HandlePostResolveDice(game *Game, player int, dice []int) {
	tally := helper.IntTally(dice)
	if helper.IntMin(tally[Die1], tally[Die2], tally[Die3]) > 0 {
		game.Boards[player].ModifyVP(2)
	}
}
