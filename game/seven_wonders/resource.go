package seven_wonders

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
