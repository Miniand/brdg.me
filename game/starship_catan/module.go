package starship_catan

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/helper"
)

const (
	ModuleLogistics = iota
	ModuleCommand
	ModuleSensor
	ModuleTrade
	ModuleScience
	ModuleProduction
)

var Modules = []int{
	ModuleLogistics,
	ModuleCommand,
	ModuleSensor,
	ModuleTrade,
	ModuleScience,
	ModuleProduction,
}

var ModuleNames = map[int]string{
	ModuleLogistics:  "logistics",
	ModuleCommand:    "command",
	ModuleSensor:     "sensor",
	ModuleTrade:      "trade",
	ModuleScience:    "science",
	ModuleProduction: "production",
}

var ModuleSummaries = map[int]string{
	ModuleLogistics: "store extra goods (2, 3, 4)",
	ModuleCommand:   "take extra actions (2, 3, 4)",
	ModuleSensor:    "peek at sector cards (0, 2, 3)",
	ModuleTrade: fmt.Sprintf(
		"buy resources from opponent for %s (0, 1, 2)",
		RenderMoney(2),
	),
	ModuleScience:    "produce science (0, 1, 2)",
	ModuleProduction: "produce trade (0, 1, 2)",
}

func ModuleTransaction(level int) Transaction {
	return Transaction{
		ResourceOre:    -1,
		ResourceCarbon: -1,
		ResourceFood:   -level,
	}
}

func ModuleDescription(module, player, level int) string {
	switch module {
	case ModuleLogistics:
		return fmt.Sprintf(
			"Store up to %d resources in each resource bay", 2+level)
	case ModuleCommand:
		return fmt.Sprintf(
			"Take up to %d actions during your flight phase", 2+level)
	case ModuleSensor:
		return fmt.Sprintf(
			"Look at the first %d sector cards of a flight, put each card on the bottom or top of the stack in any order", 1+level)
	case ModuleTrade:
		return fmt.Sprintf(
			"Buy %d resource(s) of your choice from your opponent for 2 Astro each", level)
	case ModuleScience:
		return fmt.Sprintf(
			"Produce a science point on a roll of a %s",
			strings.Join(Itoas(ScienceModuleDice(level, player)), " or "))
	case ModuleProduction:
		return fmt.Sprintf(
			"Produce a trade good on a roll of a %s",
			strings.Join(Itoas(TradeModuleDice(level, player)), " or "))
	}
	return ""
}

func ParseModule(input string) (int, error) {
	return helper.MatchStringInStringMap(input, ModuleNames)
}

func ScienceModuleDice(level, player int) []int {
	switch level {
	case 0:
		return []int{}
	case 1:
		return []int{3 - player}
	default:
		return []int{2, 3}
	}
}

func TradeModuleDice(level, player int) []int {
	switch level {
	case 0:
		return []int{}
	case 1:
		return []int{2 + player}
	default:
		return []int{2, 3}
	}
}
