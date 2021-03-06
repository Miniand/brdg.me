package red7

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	rule := g.CurrentRule()

	handCells := [][]interface{}{}
	for _, c := range helper.IntReverse(SortBySuit(g.Hands[pNum])) {
		suit := c & SuitMask
		handCells = append(handCells, []interface{}{
			RenderCard(c),
			render.Colour(SuitRulesStrs[suit], SuitCol[suit]),
		})
	}

	playerCells := [][]interface{}{
		{
			render.Bold("Player"),
			render.Bold("Hand"),
			render.Bold("Pts"),
			render.Bold("Palette"),
		},
	}
	pl := len(g.Players)
	for i := 0; i < pl; i++ {
		p := (pNum + i) % pl
		pal := render.Colour("Eliminated", render.Gray)
		if !g.Eliminated[p] {
			pal = strings.Join(RenderCards(helper.IntReverse(SortBySuit(g.Palettes[p]))), " ")
		}
		playerCells = append(playerCells, []interface{}{
			g.PlayerName(p),
			render.Centred(fmt.Sprintf("{{b}}%d{{_b}}", len(g.Hands[p]))),
			render.Centred(fmt.Sprintf("{{b}}%d{{_b}}", g.PlayerPoints(p))),
			pal,
		})
	}

	colorCells := [][]interface{}{}
	for _, s := range helper.IntReverse(Suits) {
		colorCells = append(colorCells, []interface{}{
			render.Markup(SuitStr[s], SuitCol[s], true),
			render.Colour(SuitRulesStrs[s], SuitCol[s]),
		})
	}

	rows := []interface{}{
		fmt.Sprintf("First to {{b}}%d{{_b}} points", EndPoints(len(g.Players))),
		"",
		render.Bold("Current rule"),
		render.Markup(SuitRulesStrs[rule], SuitCol[rule], true),
		"",
		render.Bold("Your hand"),
		render.Table(handCells, 0, 2),
		"",
		fmt.Sprintf("{{b}}Deck remaining:{{_b}} %d", len(g.Deck)),
		"",
		render.Table(playerCells, 0, 2),
		"",
		"",
		render.Table(colorCells, 0, 2),
	}
	return render.CentreLayout(rows, 0), nil
}

func RenderCard(card int) string {
	suit, rank := CardValues(card)
	return render.Markup(fmt.Sprintf(
		"%s%d",
		SuitAbbr[suit],
		RankVal[rank],
	), SuitCol[suit], true)
}

func RenderCards(cards []int) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = RenderCard(c)
	}
	return out
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
