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

func (c BuildCommand) Name() string { return "build" }

func (c BuildCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) < 2 {
		return "", errors.New(
			"you must pass an amount to build and the name of a thing to build")
	}

	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("you must specify an amount")
	}

	stringMap := map[int]string{
		-1: "city",
	}
	if g.Boards[pNum].Developments[DevelopmentShipping] {
		stringMap[-2] = "ship"
	}
	for _, m := range Monuments {
		stringMap[m] = MonumentValues[m].Name
	}
	thing, err := helper.MatchStringInStringMap(
		strings.Join(args[1:], " "),
		stringMap,
	)
	if err != nil {
		return "", err
	}

	switch thing {
	case -1:
		return "", g.BuildCity(pNum, amount)
	case -2:
		return "", g.BuildShip(pNum, amount)
	default:
		return "", g.BuildMonument(pNum, thing, amount)
	}
}

func (c BuildCommand) Usage(player string, context interface{}) string {
	return "{{b}}build # (thing){{_b}} to build monuments or cities using workers, or ships using cloth and wood. Eg. {{b}}build 2 great{{_b}} or {{b}}build 3 city{{_b}} or {{b}}build 1 ship{{_b}}"
}

func (g *Game) CanBuild(player int) bool {
	return g.CanBuildBuilding(player) || g.CanBuildShip(player) || g.CanTrade(player)
}

func (g *Game) CanBuildBuilding(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseBuild &&
		g.RemainingWorkers > 0
}

func (g *Game) CanBuildShip(player int) bool {
	b := g.Boards[player]
	return g.CurrentPlayer == player && g.Phase == PhaseBuild &&
		b.Developments[DevelopmentShipping] &&
		b.Goods[GoodWood] > 0 && b.Goods[GoodCloth] > 0
}

func (g *Game) BuildCity(player, amount int) error {
	if !g.CanBuildBuilding(player) {
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
	if !g.CanBuild(player) {
		g.NextPhase()
	}
	return nil
}

func (g *Game) BuildShip(player, amount int) error {
	if !g.CanBuildShip(player) {
		return errors.New("you can't build a ship at the moment")
	}
	if amount < 1 {
		return errors.New("amount must be a positive number")
	}
	if w := g.Boards[player].Goods[GoodWood]; amount > w {
		return fmt.Errorf("you only have %d wood left", w)
	}
	if c := g.Boards[player].Goods[GoodWood]; amount > c {
		return fmt.Errorf("you only have %d cloth left", c)
	}
	if g.Boards[player].Ships+amount > 5 {
		return errors.New("you can only have 5 ships")
	}

	g.Boards[player].Ships += amount
	g.Boards[player].Goods[GoodWood] -= amount
	g.Boards[player].Goods[GoodCloth] -= amount

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s built {{b}}%d ships{{_b}}",
		g.RenderName(player),
		amount,
	)))
	if !g.CanBuild(player) {
		g.NextPhase()
	}
	return nil
}

func (g *Game) BuildMonument(player, monument, amount int) error {
	if !g.CanBuildBuilding(player) {
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
	if !g.CanBuild(player) {
		g.NextPhase()
	}
	return nil
}
