package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
)

type CardPoisonQuills struct{}

func (c CardPoisonQuills) Name() string {
	return "Poison Quills"
}

func (c CardPoisonQuills) Description() string {
	twoDie := RenderDie(Die2)
	return fmt.Sprintf(
		"When you score %s%s%s, {{b}}your attack deals 2 extra damage.{{_b}}",
		twoDie, twoDie, twoDie,
	)
}

func (c CardPoisonQuills) Cost() int {
	return 3
}

func (c CardPoisonQuills) Kind() int {
	return CardKindKeep
}

func (c CardPoisonQuills) ModifyAttack(game *Game, player, damage int, attacked []int) (int, []int) {
	tally := helper.IntTally(game.CurrentRoll)
	if damage > 0 && tally[Die2] >= 3 {
		damage += 2
	}
	return damage, attacked
}
