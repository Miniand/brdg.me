package jaipur

import "github.com/Miniand/brdg.me/render"

const (
	GoodDiamond = iota
	GoodGold
	GoodSilver
	GoodCloth
	GoodSpice
	GoodLeather
	GoodCamel

	CamelBonusPoints = 5

	MinTradeBonus = 3
	MaxTradeBonus = 5
)

var TradeGoods = []int{
	GoodDiamond,
	GoodGold,
	GoodSilver,
	GoodCloth,
	GoodSpice,
	GoodLeather,
}

var Goods = append(append([]int{}, TradeGoods...), GoodCamel)

var TradeGoodTokens = map[int][]int{
	GoodDiamond: {7, 7, 5, 5, 5},
	GoodGold:    {6, 6, 5, 5, 5},
	GoodSilver:  {5, 5, 5, 5, 5},
	GoodCloth:   {5, 3, 3, 2, 2, 1, 1},
	GoodSpice:   {5, 3, 3, 2, 2, 1, 1},
	GoodLeather: {4, 3, 2, 1, 1, 1, 1, 1, 1},
}

var TradeBonuses = map[int][]int{
	5: {10, 10, 9, 8, 8},
	4: {6, 6, 5, 5, 4, 4},
	3: {3, 3, 2, 2, 2, 1, 1},
}

var GoodStrings = map[int]string{
	GoodDiamond: "diamond",
	GoodGold:    "gold",
	GoodSilver:  "silver",
	GoodCloth:   "cloth",
	GoodSpice:   "spice",
	GoodLeather: "leather",
	GoodCamel:   "camel",
}

var GoodStringsPl = map[int]string{
	GoodDiamond: "diamonds",
	GoodGold:    "golds",
	GoodSilver:  "silvers",
	GoodCloth:   "cloths",
	GoodSpice:   "spices",
	GoodLeather: "leathers",
	GoodCamel:   "camels",
}

var GoodColours = map[int]string{
	GoodDiamond: render.Red,
	GoodGold:    render.Yellow,
	GoodSilver:  render.Gray,
	GoodCloth:   render.Magenta,
	GoodSpice:   render.Green,
	GoodLeather: render.Blue,
	GoodCamel:   render.Black,
}

var CardCounts = map[int]int{
	GoodDiamond: 6,
	GoodGold:    6,
	GoodSilver:  6,
	GoodCloth:   8,
	GoodSpice:   8,
	GoodLeather: 10,
	GoodCamel:   8, // Actually 11 in game but start with 3 on the board
}

var GoodMinSales = map[int]int{
	GoodDiamond: 2,
	GoodGold:    2,
	GoodSilver:  2,
	GoodCloth:   1,
	GoodSpice:   1,
	GoodLeather: 1,
}

func Deck() []int {
	cards := []int{}
	for _, g := range Goods {
		for i := 0; i < CardCounts[g]; i++ {
			cards = append(cards, g)
		}
	}
	return cards
}
