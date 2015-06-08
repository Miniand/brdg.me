package alhambra

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type SpendCommand struct{}

func (c SpendCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("spend", 1, -1, input)
}

func (c SpendCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanSpend(pNum)
}

func (c SpendCommand) Call(player string, context interface{},
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
	return "", g.Spend(pNum, cards)
}

func (c SpendCommand) Usage(player string, context interface{}) string {
	return "{{b}}spend ## (##){{_b}} to spend cards of a single currency to buy a tile, eg. {{b}}spend r3 r4{{_b}}"
}

func (g *Game) CanSpend(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseAction
}

func (g *Game) Spend(player int, cards card.Deck) error {
	if !g.CanSpend(player) {
		return errors.New("unable to spend right now")
	}
	if cards.Len() == 0 {
		return errors.New("must specify at least one card to spend")
	}
	currency := 0
	total := 0
	counts := map[card.Card]int{}
	for i, c := range cards {
		counts[c]++
		crd := c.(Card)
		total += crd.Value
		if i == 0 {
			currency = crd.Currency
		} else {
			if crd.Currency != currency {
				return errors.New("you can only spend using cards of the same currency")
			}
		}
	}
	for c, n := range counts {
		if g.Boards[player].Cards.Contains(c) < n {
			return fmt.Errorf("you don't have enough %s in your hand", c)
		}
	}
	tile := g.Tiles[currency]
	if tile.Type == TileTypeEmpty {
		return errors.New("there isn't a tile for that currency at the moment")
	}
	if tile.Cost > total {
		return fmt.Errorf(
			"you require %d to spend that card",
			tile.Cost,
		)
	}

	// Pay for tile and add it to the tiles left to place
	for _, c := range cards {
		g.Boards[player].Cards, _ = g.Boards[player].Cards.Remove(c, 1)
	}
	g.DiscardPile = g.DiscardPile.PushMany(cards)
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
