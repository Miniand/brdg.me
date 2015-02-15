package seven_wonders

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/cost"
)

type BuildCommand struct{}

func (c BuildCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("build", 1, input)
}

func (c BuildCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanBuild(pNum)
}

func (c BuildCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which numbered card to build")
	}
	cardNum, err := strconv.Atoi(a[0])
	if err != nil || cardNum < 1 || cardNum > len(g.Hands[pNum]) {
		return "", errors.New("that is not a valid card number")
	}
	return "", g.Build(pNum, cardNum-1)
}

func (c BuildCommand) Usage(player string, context interface{}) string {
	return "{{b}}build #{{_b}} to build a card, paying the cost or getting it for free if you have a prerequisite card, eg. {{b}}build 3{{_b}}"
}

func (g *Game) CanBuild(player int) bool {
	return g.CanAction(player)
}

func (g *Game) CanBuildCard(player int, carder Carder) (
	can bool, coins []map[int]int) {
	coins = []map[int]int{}
	c := carder.GetCard()
	// See if you already have it, which means you can't get it again.
	for _, pc := range g.Cards[player] {
		if pc.(Carder).GetCard().Name == c.Name {
			can = false
			return
		}
	}
	// See if you can get it for free using another card.
	for _, freeWith := range c.FreeWith {
		if g.Cards[player].Contains(Cards[freeWith]) > 0 {
			can = true
			return
		}
	}
	return g.CanAfford(player, c.Cost)
}

func (g *Game) CanAfford(player int, c cost.Cost) (can bool, coins []map[int]int) {
	// If it costs more coin than we have, bad luck.
	if c[GoodCoin] > g.Coins[player] {
		return
	}
	// Find out if we get goods cheaper from neighbours.
	discounts := map[int]map[int]bool{
		DirLeft:  {},
		DirRight: {},
	}
	for _, pc := range g.Cards[player] {
		if disc, ok := pc.(TradeDiscounter); ok {
			dirs, goods := disc.TradeDiscount()
			for _, d := range dirs {
				for _, g := range goods {
					discounts[d][g] = true
				}
			}
		}
	}
	// Do a perm check on affordability, first favouring left, then favouring
	// right.
	for _, mode := range []struct {
		Favour int
		Dirs   []int
	}{
		{DirLeft, []int{DirLeft, DirRight}},
		{DirRight, []int{DirRight, DirLeft}},
	} {
		with := [][]cost.Cost{}
		ownerMap := map[int]int{}
		costMap := map[int]int{}
		i := 0
		// Own producing cards
		for _, pc := range g.Cards[player] {
			if prod, ok := pc.(GoodsProducer); ok {
				with = append(with, prod.GoodsProduced())
				ownerMap[i] = DirDown
				costMap[i] = 0
				i++
			}
		}
		// Neighbour trading cards, first figure out what is discounted and
		// what is full price.
		discounted := map[int][][]cost.Cost{
			DirLeft:  {},
			DirRight: {},
		}
		fullPrice := map[int][][]cost.Cost{
			DirLeft:  {},
			DirRight: {},
		}
		for _, dir := range mode.Dirs {
			for _, pc := range g.Cards[g.NumFromPlayer(player, dir)] {
				if trade, ok := pc.(GoodsTrader); ok {
					goods := trade.GoodsTraded()
					for _, perm := range goods {
						for g := range perm {
							switch discounts[dir][g] {
							case true:
								discounted[dir] = append(discounted[dir], goods)
							case false:
								fullPrice[dir] = append(fullPrice[dir], goods)
							}
							break
						}
					}
				}
			}
		}
		// Add all the discounted goods to the with slice to prioritise them.
		for _, dir := range mode.Dirs {
			for _, g := range discounted[dir] {
				with = append(with, g)
				ownerMap[i] = dir
				costMap[i] = 1
				i++
			}
		}
		// Add all the full price goods last after discounted ones.
		for _, dir := range mode.Dirs {
			for _, g := range fullPrice[dir] {
				with = append(with, g)
				ownerMap[i] = dir
				costMap[i] = 2
				i++
			}
		}
		// Check if we can afford them.
		canAfford, canWith := cost.CanAffordPerm(c.Drop(GoodCoin), with)
		if !canAfford {
			return
		}
		// Find the cheapest alternative.
		minSumCoins := 0
		maxToPriority := 0
		first := true
		minCoins := map[int]int{}
		for _, w := range canWith {
			curCoins := map[int]int{}
			sumCoins := 0
			for i, wc := range w {
				owner := ownerMap[i]
				if owner == DirDown {
					continue
				}
				amt := wc.Sum() * costMap[i]
				sumCoins += amt
				curCoins[owner] += amt
			}
			curCoins = TrimIntMap(curCoins)
			if first || sumCoins < minSumCoins ||
				sumCoins == minSumCoins && curCoins[mode.Favour] > maxToPriority {
				minSumCoins = sumCoins
				maxToPriority = curCoins[mode.Favour]
				minCoins = curCoins
			}
			first = false
		}
		if minSumCoins == 0 {
			// We can afford it by ourselves.
			can = true
			return
		}
		if minSumCoins > g.Coins[player]-c[GoodCoin] {
			// Can't afford that many coins.
			can = false
			return
		}
		unique := true
		for _, coin := range coins {
			if reflect.DeepEqual(minCoins, coin) {
				unique = false
				break
			}
		}
		if unique {
			coins = append(coins, minCoins)
		}
	}
	can = true
	return
}

func (g *Game) Build(player, cardNum int) error {
	crd := g.Hands[player][cardNum].(Carder)
	can, coins := g.CanBuildCard(player, crd)
	if !can {
		return errors.New("cannot build that card")
	}
	action := &BuildAction{
		Card: cardNum,
	}
	if len(coins) <= 1 {
		action.Chosen = true
	}
	g.Actions[player] = action
	g.CheckHandComplete()
	return nil
}
