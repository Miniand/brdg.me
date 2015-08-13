package jaipur

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type TakeCommand struct{}

func (c TakeCommand) Name() string { return "take" }

func (c TakeCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("please specify what to take")
	}
	takeGoodStrs := []string{}
	forGoodStrs := []string{}
	afterFor := false
	for _, a := range args {
		if !afterFor && strings.ToLower(a) == "for" {
			afterFor = true
			continue
		}
		if afterFor {
			forGoodStrs = append(forGoodStrs, a)
		} else {
			takeGoodStrs = append(takeGoodStrs, a)
		}
	}
	takeGoods := []int{}
	forGoods := []int{}
	if len(takeGoodStrs) > 0 {
		takeGoods, err = ParseGoodStr(strings.Join(takeGoodStrs, " "))
		if err != nil {
			return "", err
		}
	}
	if len(forGoodStrs) > 0 {
		forGoods, err = ParseGoodStr(strings.Join(forGoodStrs, " "))
		if err != nil {
			return "", err
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

	forString := ""
	if len(forGoods) > 0 {
		forString = fmt.Sprintf(
			" for %s",
			render.CommaList(RenderGoods(forGoods)),
		)
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took %s%s",
		g.RenderName(player),
		render.CommaList(RenderGoods(takeGoods)),
		forString,
	)))

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
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took {{b}}%d %s{{_b}} from the market",
		g.RenderName(player),
		numCamels,
		render.Colour(helper.Plural(numCamels, "camel"), GoodColours[GoodCamel]),
	)))
	g.Camels[player] += numCamels
	g.Market = helper.IntRemove(GoodCamel, g.Market, -1)
	if g.ReplenishMarket() {
		g.NextPlayer()
	}
	return nil
}
