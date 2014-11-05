package king_of_tokyo

import "fmt"

type CardMonsterBatteries struct{}

func (c CardMonsterBatteries) Name() string {
	return "Monster Batteries"
}

func (c CardMonsterBatteries) Description() string {
	return fmt.Sprintf(
		"When you purchase Monster Batteries, put as many %s as you want on it from your reserve. Match this from the bank. At the start of each turn {{b}}take %s off and add them to your reserve.{{_b}}  When there are no %s left discard this card.",
		EnergySymbol,
		RenderEnergy(2),
		EnergySymbol,
	)
}

func (c CardMonsterBatteries) Cost() int {
	return 2
}

func (c CardMonsterBatteries) Kind() int {
	return CardKindKeep
}
