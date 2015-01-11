package king_of_tokyo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type BuyCommand struct{}

func (c BuyCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("buy", 1, -1, input)
}

func (c BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanBuy(pNum)
}

func (c BuyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) == 0 {
		return "", errors.New("you must specify which card you want to buy")
	}
	cardNames := []string{}
	for _, c := range g.Buyable(pNum) {
		cardNames = append(cardNames, c.Name())
	}
	cNum, err := helper.MatchStringInStrings(strings.Join(a, " "), cardNames)
	if err != nil {
		return "", err
	}
	return "", g.Buy(pNum, cNum)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy (name){{_b}} to buy a card, eg {{b}}buy friend{{_b}} to buy the \"Friend of Children\" card"
}

func (g *Game) CanBuy(player int) bool {
	if g.CurrentPlayer != player {
		return false
	}
	return g.Phase == PhaseBuy && len(g.Buyable(player)) > 0
}

func (g *Game) Buy(player, cardNum int) error {
	if !g.CanBuy(player) {
		return errors.New("you can't buy at the moment")
	}
	if cardNum < 0 {
		return errors.New("the card number must be positive")
	}
	buyable := g.Buyable(player)
	if l := len(buyable); cardNum >= l {
		return fmt.Errorf("the card number must be less than %d", l)
	}
	things := g.Boards[player].Things()
	c := buyable[cardNum]
	cost := c.Cost()
	for _, t := range things {
		if costMod, ok := t.(CardCostModifier); ok {
			cost = costMod.ModifyCardCost(g, player, cost)
		}
	}
	if g.Boards[player].Energy < cost {
		return fmt.Errorf(
			"you require %s to buy that card",
			RenderEnergy(cost),
		)
	}
	if c.Kind() == CardKindKeep {
		g.Boards[player].Cards = append(g.Boards[player].Cards, c)
	} else {
		g.Discard = append(g.Discard, c)
	}
	g.Boards[player].Energy -= cost
	if cardNum < len(g.FaceUpCards) {
		g.FaceUpCards = append(g.FaceUpCards[:cardNum], g.FaceUpCards[cardNum+1:]...)
	} else {
	}
	if len(g.Deck) > 0 {
		g.FaceUpCards = append(g.FaceUpCards, g.Deck[0])
		g.Deck = g.Deck[1:]
	}
	if postBuy, ok := c.(PostCardBuyHandler); ok {
		postBuy.HandlePostCardBuy(g, player, c, cost)
	}
	for _, t := range things {
		if postBuy, ok := t.(PostCardBuyHandler); ok {
			postBuy.HandlePostCardBuy(g, player, c, cost)
		}
	}
	return nil
}

type BuyableCard struct {
	Card     CardBase
	FromText string
}

func (g *Game) Buyable(player int) []BuyableCard {
	cards := []BuyableCard{}
	for _, c := range g.FaceUpCards {
		cards = append(cards, BuyableCard{
			Card: c,
		})
	}
	return cards
}
