package starship_catan

import "github.com/Miniand/brdg.me/game/helper"

const (
	ResourceAny = iota
	ResourceFood
	ResourceFuel
	ResourceCarbon
	ResourceOre
	ResourceScience
	ResourceTrade
	ResourceAstro
	ResourceColonyShip
	ResourceTradeShip
	ResourceBooster
	ResourceCannon
)

var Resources = []int{
	ResourceAny,
	ResourceFood,
	ResourceFuel,
	ResourceCarbon,
	ResourceOre,
	ResourceScience,
	ResourceTrade,
	ResourceAstro,
	ResourceColonyShip,
	ResourceTradeShip,
	ResourceBooster,
	ResourceCannon,
}

var ResourceNames = map[int]string{
	ResourceAny:        "any resource",
	ResourceFood:       "food",
	ResourceFuel:       "fuel",
	ResourceCarbon:     "carbon",
	ResourceOre:        "ore",
	ResourceScience:    "science",
	ResourceTrade:      "trade",
	ResourceAstro:      "astro",
	ResourceColonyShip: "colony ship",
	ResourceTradeShip:  "trade ship",
	ResourceBooster:    "booster",
	ResourceCannon:     "cannon",
}

var ResourceColours = map[int]string{
	ResourceAny:        "green",
	ResourceFood:       "red",
	ResourceFuel:       "gray",
	ResourceCarbon:     "cyan",
	ResourceOre:        "black",
	ResourceScience:    "magenta",
	ResourceTrade:      "yellow",
	ResourceAstro:      "green",
	ResourceColonyShip: "cyan",
	ResourceTradeShip:  "yellow",
	ResourceBooster:    "red",
	ResourceCannon:     "blue",
}

var ColonyShipTransaction = Transaction{
	ResourceOre:        -1,
	ResourceFuel:       -1,
	ResourceFood:       -1,
	ResourceColonyShip: 1,
}

var TradeShipTransaction = Transaction{
	ResourceOre:       -1,
	ResourceFuel:      -1,
	ResourceTrade:     -1,
	ResourceTradeShip: 1,
}

var Goods = []int{
	ResourceFood,
	ResourceFuel,
	ResourceCarbon,
	ResourceOre,
	ResourceTrade,
}

var Buildables = []int{
	ResourceTradeShip,
	ResourceColonyShip,
	ResourceBooster,
	ResourceCannon,
}

func ResourceNameArr(resources []int) []string {
	names := make([]string, len(resources))
	for i, r := range resources {
		names[i] = ResourceNames[r]
	}
	return names
}

func ResourceNameMap(resources []int) map[int]string {
	rm := map[int]string{}
	for _, r := range resources {
		rm[r] = ResourceNames[r]
	}
	return rm
}

func ParseBuildable(input string) (int, error) {
	return helper.MatchStringInStringMap(input, ResourceNameMap(Buildables))
}

func ParseGood(input string) (int, error) {
	return helper.MatchStringInStringMap(input, ResourceNameMap(Goods))
}

func ParseResource(input string) (int, error) {
	return helper.MatchStringInStringMap(input, ResourceNames)
}
