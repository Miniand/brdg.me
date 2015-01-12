package splendor

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

var ResourceColours = map[int]string{
	Diamond:  render.Gray,
	Sapphire: render.Blue,
	Emerald:  render.Green,
	Ruby:     render.Red,
	Onyx:     render.Black,
	Gold:     render.Yellow,
	Prestige: render.Magenta,
}

var ResourceStrings = map[int]string{
	Diamond:  "diamond",
	Sapphire: "sapphire",
	Emerald:  "emerald",
	Ruby:     "ruby",
	Onyx:     "onyx",
	Gold:     "gold",
	Prestige: "prestige",
}

var ResourceAbbr = map[int]string{
	Diamond:  "Di",
	Sapphire: "Sa",
	Emerald:  "Em",
	Ruby:     "Ru",
	Onyx:     "On",
	Gold:     "Go",
	Prestige: "Pr",
}

func splitCards(cards []Card, n int) [][]Card {
	rows := [][]Card{}
	for i := 0; i < len(cards)/n; i++ {
		rows = append(rows, cards[n*i:n*(i+1)])
	}
	return rows
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	// Board
	longestRow := 0
	for _, r := range g.Board {
		if l := len(r); l > longestRow {
			longestRow = l
		}
	}
	header := []interface{}{""}
	for i := 0; i < longestRow; i++ {
		header = append(header, render.Centred(
			render.Colour(fmt.Sprintf("%c", 'A'+i), render.Gray)))
	}
	table := [][]interface{}{header}
	for l, r := range g.Board {
		upper := []interface{}{
			render.Colour(fmt.Sprintf("Level %d", l+1), render.Gray),
		}
		lower := []interface{}{
			"",
		}
		for _, c := range r {
			canAfford := g.PlayerBoards[pNum].CanAfford(c.Cost)
			upper = append(upper, render.Centred(
				render.Markup(RenderCardBonusVP(c), "", canAfford)))
			lower = append(lower, render.Centred(
				render.Markup(RenderCardCost(c), "", canAfford)))
		}
		table = append(table, upper, lower, []interface{}{})
	}
	return render.Table(table, 0, 3), nil
}

func RenderResourceColour(v interface{}, r int) string {
	return render.Colour(v, ResourceColours[r])
}

func RenderCardBonusVP(c Card) string {
	parts := []string{
		RenderResourceColour(ResourceAbbr[c.Resource], c.Resource),
	}
	if c.Prestige > 0 {
		parts = append(parts, RenderResourceColour(c.Prestige, Prestige))
	}
	return strings.Join(parts, " ")
}

func RenderCardCost(c Card) string {
	parts := []string{}
	for _, r := range Resources {
		if c.Cost[r] > 0 {
			parts = append(parts, RenderResourceColour(c.Cost[r], r))
		}
	}
	return strings.Join(parts, "")
}
