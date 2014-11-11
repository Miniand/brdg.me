package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type BuildCommand struct{}

func (c BuildCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("build", 1, -1, input)
}

func (c BuildCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanBuild(pNum)
}

func (c BuildCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 2 {
		return "", errors.New(
			"you must pass an amount to build and the name of a thing to build")
	}

	amount, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("you must specify an amount")
	}

	stringMap := map[int]string{
		-1: "city",
	}
	for _, m := range Monuments {
		stringMap[m] = MonumentValues[m].Name
	}
	thing, err := helper.MatchStringInStringMap(
		strings.Join(a[1:], " "),
		stringMap,
	)
	if err != nil {
		return "", err
	}

	if thing < 0 {
		return "", g.BuildCity(pNum, amount)
	}
	return "", g.BuildMonument(pNum, thing, amount)
}

func (c BuildCommand) Usage(player string, context interface{}) string {
	return "{{b}}build # (thing){{_b}} to build monuments or cities, eg. {{b}}build 2 great{{_b}} or {{b}}build 3 city{{_b}}"
}

func (g *Game) CanBuild(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseBuild &&
		g.RemainingWorkers > 0
}

func (g *Game) BuildCity(player, amount int) error {
	if !g.CanBuild(player) {
		return errors.New("you can't build at the moment")
	}
	if amount < 1 {
		return errors.New("amount must be a positive number")
	}
	if amount > g.RemainingWorkers {
		return fmt.Errorf("you only have %d workers left", g.RemainingWorkers)
	}
	if g.Boards[player].CityProgress+amount > MaxCityProgress {
		return errors.New("that is more than what remains to be built")
	}
	initialCities := g.Boards[player].Cities()
	g.RemainingWorkers -= amount
	g.Boards[player].CityProgress += amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s used {{b}}%d{{_b}} workers on {{b}}cities{{_b}}",
		g.RenderName(player),
		amount,
	)))
	newCities := g.Boards[player].Cities()
	if newCities > initialCities {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s now has {{b}}%d cities{{_b}}",
			g.RenderName(player),
			newCities,
		)))
	}
	if g.RemainingWorkers == 0 {
		g.BuyPhase()
	}
	return nil
}

func (g *Game) BuildMonument(player, monument, amount int) error {
	if !g.CanBuild(player) {
		return errors.New("you can't build at the moment")
	}
	if amount < 1 {
		return errors.New("amount must be a positive number")
	}
	if !ContainsInt(monument, Monuments) {
		return errors.New("that isn't a valid monument")
	}
	if amount > g.RemainingWorkers {
		return fmt.Errorf("you only have %d workers left", g.RemainingWorkers)
	}
	mv := MonumentValues[monument]
	if g.Boards[player].Monuments[monument]+amount > mv.Size {
		return errors.New("that is more than what remains to be built")
	}
	g.RemainingWorkers -= amount
	g.Boards[player].Monuments[monument] += amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s used {{b}}%d{{_b}} workers on the {{b}}%s{{_b}}",
		g.RenderName(player),
		amount,
		mv.Name,
	)))
	if g.Boards[player].Monuments[monument] >= mv.Size {
		first := true
		for pNum, _ := range g.Players {
			if g.Boards[pNum].MonumentBuiltFirst[monument] {
				first = false
				break
			}
		}
		if first {
			g.Boards[player].MonumentBuiltFirst[monument] = true
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s completed the {{b}}%s{{_b}}",
			g.RenderName(player),
			mv.Name,
		)))
		g.CheckGameEndTriggered(player)
	}
	if g.RemainingWorkers == 0 {
		g.BuyPhase()
	}
	return nil
}
