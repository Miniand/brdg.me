package seven_wonders

import "github.com/Miniand/brdg.me/game/cost"

const (
	GoodCoin = iota
	GoodWood
	GoodStone
	GoodOre
	GoodClay
	GoodPapyrus
	GoodTextile
	GoodGlass

	CardKindRaw
	CardKindManufactured
	CardKindCivilian
	CardKindScientific
	CardKindCommercial
	CardKindMilitary
	CardKindGuild

	FieldMathematics
	FieldEngineering
	FieldTheology

	AttackStrength
	TokenVictory
	TokenDefeat

	WonderStage

	VP
)

var RawGoods = []int{
	GoodWood,
	GoodStone,
	GoodOre,
	GoodClay,
}

var ManufacturedGoods = []int{
	GoodPapyrus,
	GoodTextile,
	GoodGlass,
}

func SliceToCost(ints []int) []cost.Cost {
	l := len(ints)
	c := make([]cost.Cost, l)
	for i := 0; i < l; i++ {
		c[i] = cost.Cost{
			ints[i]: 1,
		}
	}
	return c
}
