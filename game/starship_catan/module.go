package starship_catan

import (
	"fmt"
	"strconv"
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
	ModuleLogistics:  "Logistics",
	ModuleCommand:    "Command",
	ModuleSensor:     "Sensor",
	ModuleTrade:      "Trade",
	ModuleScience:    "Science",
	ModuleProduction: "Production",
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
		diceStr := "2 or 3"
		if level == 1 {
			diceStr = strconv.Itoa(3 - level)
		}
		return fmt.Sprintf(
			"Produce a science point on a roll of a %s", diceStr)
	case ModuleProduction:
		diceStr := "2 or 3"
		if level == 1 {
			diceStr = strconv.Itoa(2 + level)
		}
		return fmt.Sprintf(
			"Produce a trade good on a roll of a %s", diceStr)
	}
	return ""
}
