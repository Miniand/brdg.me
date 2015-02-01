package alhambra

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type BuyCommand struct{}

func (c BuyCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("buy", 1, -1, input)
}

func (c BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanBuy(pNum)
}

func (c BuyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", ErrCouldNotFindPlayer
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which cards to use")
	}
	cards := card.Deck{}
	for _, rawCard := range a {
		c, err := ParseCard(rawCard)
		if err != nil {
			return "", err
		}
		cards = cards.Push(c)
	}
	return "", g.Buy(pNum, cards)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy ## (##){{_b}} to buy a tile using cards of a single currency, eg. {{b}}buy r3 r4{{_b}}"
}

func (g *Game) CanBuy(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseAction
}

func (g *Game) Buy(player int, cards card.Deck) error {
	if !g.CanBuy(player) {
		return errors.New("unable to buy right now")
	}
	if cards.Len() == 0 {
		return errors.New("must specify at least one card to buy with")
	}
	currency := 0
	total := 0
	for i, c := range cards {
		if g.Boards[player].Cards.Contains(c) == 0 {
			return fmt.Errorf("you don't have %s in your hand", c)
		}
		crd := c.(Card)
		total += crd.Value
		if i == 0 {
			currency = crd.Currency
		} else {
			if crd.Currency != currency {
				return errors.New("you can only buy using cards of the same currency")
			}
		}
	}
	tile := g.Tiles[currency]
	if tile.Type == TileTypeEmpty {
		return errors.New("there isn't a tile for that currency at the moment")
	}
	if tile.Cost > total {
		return fmt.Errorf(
			"you require %d to buy that card",
			tile.Cost,
		)
	}

	// Pay for tile and add it to the tiles left to place
	for _, c := range cards {
		g.Boards[player].Cards, _ = g.Boards[player].Cards.Remove(c, 1)
	}
	g.Boards[player].Place = append(g.Boards[player].Place, tile)
	g.Tiles[currency] = Tile{}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s bought %s (cost %d) using %s",
		g.PlayerName(player),
		RenderTileAbbr(tile.Type),
		tile.Cost,
		RenderCards(cards),
	)))

	if total != tile.Cost {
		g.NextPhase()
	}
	return nil
}
