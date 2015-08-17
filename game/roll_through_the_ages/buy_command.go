package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type BuyCommand struct{}

func (c BuyCommand) Name() string { return "buy" }

func (c BuyCommand) Call(
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
	if err != nil || len(args) < 1 {
		return "", errors.New("you must specify which development to buy")
	}

	development, err := helper.MatchStringInStringMap(args[0], DevelopmentNameMap())
	if err != nil {
		return "", err
	}

	goods := []int{}
	if len(args) == 2 && strings.ToLower(args[1]) == "all" {
		for good, num := range g.Boards[pNum].Goods {
			if num > 0 {
				goods = append(goods, good)
			}
		}
	} else {
		for _, input := range args[1:] {
			good, err := helper.MatchStringInStringMap(input, GoodStrings)
			if err != nil {
				return "", err
			}
			goods = append(goods, good)
		}
	}

	return "", g.BuyDevelopment(pNum, development, goods)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy (thing){{_b}} to buy developments, eg. {{b}}buy irrigation{{_b}}.  If you don't have enough coins and need to use goods, list the goods after the development name, eg. {{b}}buy irrigation wood stone{{_b}} or {{b}}buy irrigation all{{_b}} to use all your goods"
}

func (g *Game) CanBuy(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseBuy
}

func (g *Game) BuyDevelopment(player, development int, goods []int) error {
	if !g.CanBuy(player) {
		return errors.New("you can't buy at the moment")
	}
	if g.Boards[player].Developments[development] {
		return errors.New("you already have that development")
	}
	dv, ok := DevelopmentValues[development]
	if !ok {
		return errors.New("invalid development")
	}

	total := g.RemainingCoins
	usedGoods := map[int]bool{}
	for _, good := range goods {
		if usedGoods[good] {
			continue
		}
		total += GoodValue(good, g.Boards[player].Goods[good])
		usedGoods[good] = true
	}
	if total < dv.Cost {
		return fmt.Errorf(
			`you require %d but your coins and specified goods only amount to %d, you may need to add more goods`,
			dv.Cost,
			total,
		)
	}

	suffix := ""
	if len(usedGoods) > 0 {
		suffixParts := []string{}
		for good, _ := range usedGoods {
			suffixParts = append(suffixParts, RenderGoodName(good))
		}
		suffix = fmt.Sprintf(", using %s", strings.Join(suffixParts, ", "))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s bought the {{b}}%s development{{_b}}%s`,
		g.RenderName(player),
		dv.Name,
		suffix,
	)))
	g.Boards[player].Developments[development] = true
	for _, good := range goods {
		g.Boards[player].Goods[good] = 0
	}
	g.CheckGameEndTriggered(player)

	g.NextPhase()
	return nil
}
