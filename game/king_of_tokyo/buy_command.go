package king_of_tokyo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

const (
	BuyFromFaceUp = iota
	BuyFromDeck
	BuyFromPlayer
)

type BuyCommand struct{}

func (c BuyCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("buy", 1, -1, input)
}

func (c BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanBuy(pNum)
}

func (c BuyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) == 0 {
		return "", errors.New("you must specify which card you want to buy")
	}
	cNum, err := strconv.Atoi(a[0])
	if err == nil {
		// Player passed a number, make it zero based
		cNum -= 1
	} else {
		// Try to match a string
		cardNames := []string{}
		for _, c := range g.Buyable(pNum) {
			cardNames = append(cardNames, c.Card.Name())
		}
		cNum, err = helper.MatchStringInStrings(strings.Join(a, " "), cardNames)
		if err != nil {
			return "", err
		}
	}
	return "", g.Buy(pNum, cNum)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy (something){{_b}} to buy a card using the name or number of the card, eg {{b}}buy friend{{_b}} to buy \"Friend of Children\" or {{b}}buy 1{{_b}} to buy card #1"
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
	cost := c.Card.Cost()
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
	if c.Card.Kind() == CardKindKeep {
		g.Boards[player].Cards = append(g.Boards[player].Cards, c.Card)
	} else {
		g.Discard = append(g.Discard, c.Card)
	}
	g.Boards[player].Energy -= cost
	// Remove the card from where it was bought from.
	switch c.From {
	case BuyFromFaceUp:
		g.FaceUpCards = append(g.FaceUpCards[:cardNum], g.FaceUpCards[cardNum+1:]...)
		if len(g.Deck) > 0 {
			g.FaceUpCards = append(g.FaceUpCards, g.Deck[0])
			g.Deck = g.Deck[1:]
		}
	case BuyFromDeck:
		for i, dc := range g.Deck {
			if dc == c.Card {
				g.Deck = append(g.Deck[:i], g.Deck[i+1:]...)
				break
			}
		}
	}
	if postBuy, ok := c.Card.(PostCardBuyHandler); ok {
		postBuy.HandlePostCardBuy(g, player, c.Card, cost)
	}
	for _, t := range things {
		if postBuy, ok := t.(PostCardBuyHandler); ok {
			postBuy.HandlePostCardBuy(g, player, c.Card, cost)
		}
	}
	return nil
}

type BuyableCard struct {
	Card       CardBase
	From       int
	FromPlayer int
}

func (bc BuyableCard) FromString(g *Game) string {
	switch bc.From {
	case BuyFromDeck:
		return "from the deck"
	case BuyFromPlayer:
		return fmt.Sprintf("from %s", g.RenderName(bc.FromPlayer))
	default:
		return ""
	}
}

func (g *Game) Buyable(player int) []BuyableCard {
	cards := []BuyableCard{}
	for _, c := range g.FaceUpCards {
		cards = append(cards, BuyableCard{
			Card: c,
		})
	}
	for _, t := range g.Boards[player].Things() {
		if bm, ok := t.(BuyableModifier); ok {
			cards = bm.ModifyBuyable(g, player, cards)
		}
	}
	return cards
}
