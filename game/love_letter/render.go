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
	rows := []interface{}{}
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
	return render.CentreLayout(rows, 0), nil
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func RenderCard(card int) string {
	return render.Bold(fmt.Sprintf(
		"%s (%d)",
		Cards[card].Name(),
		Cards[card].Number(),
	))
}

func RenderCards(cards []int) []string {
	strs := make([]string, len(cards))
	for i, c := range cards {
		strs[i] = RenderCard(c)
	}
	return strs
}
