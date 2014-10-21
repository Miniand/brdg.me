package starship_catan

import (
	"errors"
	"strings"
)

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
	ResourceAny:     "green",
	ResourceFood:    "red",
	ResourceFuel:    "gray",
	ResourceCarbon:  "cyan",
	ResourceOre:     "black",
	ResourceScience: "magenta",
	ResourceTrade:   "yellow",
}

func ParseResource(input string) (int, error) {
	in := []byte(strings.ToLower(input))
	skipped := map[int]bool{}
	for i, b := range in {
		found := 0
		foundR := 0
		for r, rName := range ResourceNames {
			if skipped[r] || b != rName[i] {
				skipped[r] = true
				continue
			}
			found += 1
			foundR = r
		}
		switch found {
		case 0:
			break
		case 1:
			return foundR, nil
		}
	}
	return 0, errors.New("could not find a unique resource for that input")
}
