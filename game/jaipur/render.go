package jaipur

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	buf := bytes.Buffer{}
	remainingRounds := 3 - g.RoundWins[0] - g.RoundWins[1]
	leaderText := ""
	leader := -1
	if g.RoundWins[0] > g.RoundWins[1] {
		leader = 0
	} else if g.RoundWins[1] > g.RoundWins[0] {
		leader = 1
	}
	if leader != -1 {
		leaderText = fmt.Sprintf(" %s is in the lead.", g.RenderName(leader))
	}
	cells := [][]interface{}{
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"There %s {{b}}%d{{_b}} %s remaining.%s",
			helper.Plural(remainingRounds, "is"),
			remainingRounds,
			helper.Plural(remainingRounds, "round"),
			leaderText,
		)), 7}},
		{},
		{render.CellSpan{render.Centred(render.Bold("Sale prices")), 7}},
		{
			render.CellSpan{render.Centred(render.Markup(
				"Rare",
				render.Gray,
				true,
			)), 3},
			"", // Spacer
			render.CellSpan{render.Centred(render.Markup(
				"Common",
				render.Gray,
				true,
			)), 3},
		},
	}

	subHeading := []interface{}{}
	for i, good := range TradeGoods {
		if i == 3 {
			// Spacer
			subHeading = append(subHeading, "")
		}
		subHeading = append(
			subHeading,
			render.Centred(render.Markup(
				GoodStrings[good],
				GoodColours[good],
				true,
			)),
		)
	}
	cells = append(cells, subHeading)

	i := 0
	for {
		hasContent := false
		row := []interface{}{}
		for gi, good := range TradeGoods {
			if gi == 3 {
				// Spacer
				row = append(row, "")
			}
			gl := len(g.Goods[good])
			if gl > i {
				hasContent = true
				row = append(row, render.Centred(render.Markup(
					g.Goods[good][gl-i-1],
					GoodColours[good],
					true,
				)))
			} else {
				row = append(row, "")
			}
		}
		if !hasContent {
			break
		}
		cells = append(cells, row)
		i++
	}

	cells = append(cells, [][]interface{}{
		{render.CellSpan{render.Centred(render.Bold("Bonuses for selling")), 7}},
		{render.CellSpan{render.Centred(render.Table([][]interface{}{{
			fmt.Sprintf(
				"{{b}}3{{_b}}: %d left",
				len(g.Bonuses[3]),
			),
			fmt.Sprintf(
				"{{b}}4{{_b}}: %d left",
				len(g.Bonuses[4]),
			),
			fmt.Sprintf(
				"{{b}}5 or more{{_b}}: %d left",
				len(g.Bonuses[5]),
			),
		}}, 0, 4)), 7}},
	}...)

	opponentNum := (pNum + 1) % 2
	camelStr := "no"
	if g.Camels[opponentNum] > 0 {
		camelStr = "some"
	}
	cells = append(cells, [][]interface{}{
		{},
		{render.CellSpan{render.Centred(render.Bold("Market")), 7}},
		{render.CellSpan{render.Centred(
			strings.Join(RenderGoods(helper.IntSort(g.Market)), "  ")), 7}},
		{},
		{render.CellSpan{render.Centred(render.Bold("You have")), 7}},
		{render.CellSpan{render.Centred(
			strings.Join(RenderGoods(helper.IntSort(g.Hands[pNum])), "  ")), 7}},
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"%d %s",
			g.Camels[pNum],
			render.Markup(
				helper.Plural(g.Camels[pNum], GoodStrings[GoodCamel]),
				GoodColours[GoodCamel],
				true,
			),
		)), 7}},
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"%d %s",
			len(g.Tokens[pNum]),
			helper.Plural(len(g.Tokens[pNum]), "point token"),
		)), 7}},
		{},
		{render.CellSpan{render.Centred(render.Bold("Your opponent has")), 7}},
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"%d %s",
			len(g.Hands[opponentNum]),
			helper.Plural(len(g.Hands[opponentNum]), "good"),
		)), 7}},
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"%s %s",
			camelStr,
			render.Markup(
				helper.Plural(2, GoodStrings[GoodCamel]),
				GoodColours[GoodCamel],
				true,
			),
		)), 7}},
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"%d %s",
			len(g.Tokens[opponentNum]),
			helper.Plural(len(g.Tokens[opponentNum]), "point token"),
		)), 7}},
		{},
		{render.CellSpan{render.Centred(fmt.Sprintf(
			"{{b}}%d{{_b}} %s left in the deck",
			len(g.Deck),
			helper.Plural(len(g.Deck), "card"),
		)), 7}},
	}...)

	buf.WriteString(render.Table(cells, 0, 2))
	return buf.String(), nil
}

func RenderGood(good int) string {
	return render.Markup(GoodStrings[good], GoodColours[good], true)
}

func RenderGoods(goods []int) []string {
	if goods == nil {
		return nil
	}
	strs := make([]string, len(goods))
	for i, g := range goods {
		strs[i] = RenderGood(g)
	}
	return strs
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
