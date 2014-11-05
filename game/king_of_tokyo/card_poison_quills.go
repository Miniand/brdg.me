package king_of_tokyo

import "fmt"

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
