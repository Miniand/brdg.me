package jaipur

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type TakeCommand struct{}

func (c TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("take", 1, -1, input)
}

func (c TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	return found && g.CanTake(pNum)
}

func (c TakeCommand) Call(
	player string,
	context interface{},
	args []string,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	takeGoods := []int{}
	forGoods := []int{}
	afterFor := false
	for _, a := range command.ExtractNamedCommandArgs(args) {
		if !afterFor && strings.ToLower(a) == "for" {
			afterFor = true
			continue
		}
		good, err := helper.MatchStringInStringMap(a, GoodStrings)
		if err != nil {
			return "", err
		}
		if afterFor {
			forGoods = append(forGoods, good)
		} else {
			takeGoods = append(takeGoods, good)
		}
	}
	return "", g.Take(pNum, takeGoods, forGoods)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take [goods] (for [goods]){{_b}} to take cards from the market, eg. {{b}}take dia{{_b}} or {{b}}take dia silv for camel spi{{_b}}"
}

func (g *Game) CanTake(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Take(player int, takeGoods, forGoods []int) error {
	if !g.CanTake(player) {
		return errors.New("can't take at the moment")
	}
	numTake := len(takeGoods)
	numFor := len(forGoods)
	if numTake == 0 {
		return errors.New("you must specify a good to take")
	}
	if numTake == 1 && numFor != 0 {
		return errors.New("if you are taking a single good you can't put any back")
	}
	if numTake > 1 && numTake != numFor {
		return errors.New("if you are taking more than one good you must trade for the same number of goods in your hand, eg take gold dia for lea lea")
	}
	if numTake > 1 && helper.IntCount(GoodCamel, takeGoods) > 0 {
		return errors.New("the only way to take camels is using \"take camel\" which will take all the camels from the market and replace them with cards drawn from the deck")
	}
	if numTake == 1 && takeGoods[0] == GoodCamel {
		if numFor > 0 {
			return errors.New("when taking camels you don't trade for cards in your hand, they are replaced from the deck instead")
		}
		return g.TakeCamels(player)
	}
	handSizeAfter := len(g.Hands[player]) + numTake
	for _, good := range forGoods {
		if good != GoodCamel {
			handSizeAfter--
		}
	}
	if handSizeAfter > 7 {
		return errors.New("that would exceed your hand size of 7")
	}

	// Make sure we aren't trading the same type of good.
	takeMap := helper.IntTally(takeGoods)
	for _, good := range forGoods {
		if takeMap[good] > 0 {
			return errors.New("you can't trade the same type of good")
		}
	}

	// Make sure the market has what we want.
	marketMap := helper.IntTally(g.Market)
	for good, n := range takeMap {
		if marketMap[good] < n {
			return fmt.Errorf(
				"the market only has %d %s",
				marketMap[good],
				GoodStrings[good],
			)
		}
	}

	// Make sure we have enough goods to trade out of our hand.
	forMap := helper.IntTally(forGoods)
	handMap := helper.IntTally(g.Hands[player])
	handMap[GoodCamel] = g.Camels[player]
	for good, n := range forMap {
		if handMap[good] < n {
			return fmt.Errorf(
				"you only have %d %s",
				handMap[good],
				GoodStrings[good],
			)
		}
	}

	for good, n := range takeMap {
		g.Market = helper.IntRemove(good, g.Market, n)
	}
	for good, n := range forMap {
		switch good {
		case GoodCamel:
			g.Camels[player] -= n
		default:
			g.Hands[player] = helper.IntRemove(good, g.Hands[player], n)
		}
	}
	g.Market = append(g.Market, forGoods...)
	g.Hands[player] = append(g.Hands[player], takeGoods...)
	if g.ReplenishMarket() {
		g.NextPlayer()
	}
	return nil
}

func (g *Game) TakeCamels(player int) error {
	if !g.CanTake(player) {
		return errors.New("can't take at the moment")
	}
	numCamels := helper.IntCount(GoodCamel, g.Market)
	if numCamels == 0 {
		return errors.New("there are no camels in the market")
	}
	g.Camels[player] += numCamels
	g.Market = helper.IntRemove(GoodCamel, g.Market, -1)
	if g.ReplenishMarket() {
		g.NextPlayer()
	}
	return nil
}
