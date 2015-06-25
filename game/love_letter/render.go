package love_letter

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

	leader := helper.IntMax(g.Points...)
	rows := []interface{}{fmt.Sprintf(
		"The leader has {{b}}%d %s{{_b}}, the game will end at {{b}}%d points{{_b}}",
		leader,
		helper.Plural(leader, "point"),
		endScores[len(g.Players)],
	), ""}

	if g.Eliminated[pNum] {
		rows = append(
			rows,
			render.Bold("You have been eliminated from this round"),
		)
	} else {
		rows = append(
			rows,
			render.Bold(fmt.Sprintf(
				"Your %s",
				helper.Plural(len(g.Hands[pNum]), "card"),
			)),
			strings.Join(RenderCards(g.Hands[pNum]), "   "),
		)
	}

	playerTable := [][]interface{}{
		{},
		{
			render.Bold("Player"),
			render.Centred(render.Bold("Pts")),
			render.Centred(render.Bold("Status")),
			render.Centred(render.Bold("Discards")),
		},
	}
	for p := range g.Players {
		status := render.Markup("active", render.Green, true)
		if g.Eliminated[p] {
			status = render.Markup("eliminated", render.Gray, false)
		} else if g.Protected[p] {
			status = render.Markup("protected", render.Black, true)
		}
		playerTable = append(playerTable, []interface{}{
			g.RenderName(p),
			render.Centred(render.Bold(g.Points[p])),
			render.Centred(status),
			strings.Join(RenderCards(g.Discards[p]), "  "),
		})
	}
	rows = append(rows,
		render.Table(playerTable, 0, 2),
		"",
		render.Bold(fmt.Sprintf("Cards remaining: %d", len(g.Deck))),
	)

	helpTable := [][]interface{}{
		{},
		{},
		{render.Bold("Card"), render.Bold("#"), render.Bold("Description")},
	}
	for c := Princess; c >= Guard; c-- {
		helpTable = append(helpTable, []interface{}{
			RenderCard(c),
			helper.IntCount(c, Deck),
			render.Colour(Cards[c].Text(), render.Gray),
		})
	}
	rows = append(rows, render.Table(helpTable, 0, 2))

	return render.CentreLayout(rows, 0), nil
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func RenderCard(card int) string {
	return render.Markup(fmt.Sprintf(
		"%s (%d)",
		Cards[card].Name(),
		Cards[card].Number(),
	), Cards[card].Colour(), true)
}

func RenderCards(cards []int) []string {
	strs := make([]string, len(cards))
	for i, c := range cards {
		strs[i] = RenderCard(c)
	}
	return strs
}
