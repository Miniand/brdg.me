package seven_wonders

import "fmt"

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
	CardKindWonder

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

var Goods = []int{
	GoodWood,
	GoodStone,
	GoodOre,
	GoodClay,
	GoodPapyrus,
	GoodTextile,
	GoodGlass,
}

var CardKinds = []int{
	CardKindRaw,
	CardKindManufactured,
	CardKindCivilian,
	CardKindScientific,
	CardKindCommercial,
	CardKindMilitary,
	CardKindGuild,
}

var Fields = []int{
	FieldMathematics,
	FieldEngineering,
	FieldTheology,
}

func (g *Game) PlayerResourceCount(player, resource int) int {
	switch {
	case resource == GoodCoin:
		return g.Coins[player]
	case resource == TokenVictory:
		return g.VictoryTokens[player]
	case resource == TokenDefeat:
		return g.DefeatTokens[player]
	case resource == VP:
		sum := 0
		for _, c := range g.Cards[player] {
			if vp, ok := c.(VictoryPointer); ok {
				sum += vp.VictoryPoints(player, g)
			}
		}
		return sum
	case resource == CardKindWonder || InInts(resource, CardKinds):
		sum := 0
		for _, c := range g.Cards[player] {
			if c.(Carder).GetCard().Kind == resource {
				sum++
			}
		}
		return sum
	case InInts(resource, RawGoods) || InInts(resource, ManufacturedGoods):
		sum := 0
		for _, c := range g.Cards[player] {
			if prod, ok := c.(GoodsProducer); ok {
				produced := prod.GoodsProduced()
				if len(produced) != 1 {
					// We ignore ones that produce different kinds as we only
					// want a minimum.
					continue
				}
				sum += produced[0][resource]
			}
		}
		return sum
	case resource == AttackStrength:
		sum := 0
		for _, c := range g.Cards[player] {
			if attacker, ok := c.(Attacker); ok {
				sum += attacker.AttackStrength()
			}
		}
		return sum
	default:
		panic(fmt.Sprintf("Good %d not implemented", resource))
	}
}
