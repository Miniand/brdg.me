package age_of_war

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	layout := [][]interface{}{}

	if g.CurrentlyAttacking != -1 {
		layout = append(layout, [][]interface{}{
			{render.Centred(render.Bold("Currently attacking\n"))},
			{render.Centred(Castles[g.CurrentlyAttacking].Render(
				g.Conquered[g.CurrentlyAttacking],
			))},
			{},
		}...)
	}

	layout = append(layout, [][]interface{}{
		{render.Centred(render.Bold("Current roll"))},
		{render.Centred(strings.Join(RenderDice(g.CurrentRoll), "   "))},
		{},
		{},
		{render.Centred(render.Bold("Available castles"))},
		{},
	}...)
	layout = append(layout, []interface{}{
		render.Centred(RenderCastles(Castles, false)),
	})
	return render.Table(layout, 0, 0), nil
}

func RenderCastles(castles []Castle, stealing bool) string {
	cells := [][]interface{}{}
	row := []interface{}{}
	lastClan := -1
	for _, c := range castles {
		if lastClan != -1 && c.Clan != lastClan {
			cells = append(cells, []interface{}{render.Centred(render.Table(
				[][]interface{}{row}, 0, 6))})
			row = []interface{}{}
		}
		row = append(row, c.Render(stealing))
		lastClan = c.Clan
	}
	if len(row) > 0 {
		cells = append(cells, []interface{}{render.Centred(render.Table(
			[][]interface{}{row}, 0, 6))})
	}
	return render.Table(cells, 2, 6)
}

func (c Castle) Render(stealing bool) string {
	return render.Table(c.RenderCells(stealing), 0, 0)
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func RenderDie(die int) string {
	return render.Markup(DiceStrings[die], DiceColours[die], true)
}

func RenderDice(dice []int) []string {
	l := len(dice)
	if l == 0 {
		return []string{}
	}
	strs := make([]string, l)
	for i, d := range dice {
		strs[i] = RenderDie(d)
	}
	return strs
}

func RenderInf(n int) string {
	return render.Markup(fmt.Sprintf("%d inf", n), InfantryColour, true)
}

func (c Castle) RenderName() string {
	return render.Markup(c.Name, ClanColours[c.Clan], true)
}

func RenderClan(clan int) string {
	return render.Markup(ClanNames[clan], ClanColours[clan], true)
}

func (c Castle) RenderCells(stealing bool) [][]interface{} {
	cells := [][]interface{}{
		{render.Centred(fmt.Sprintf(
			"%s (%d)",
			c.RenderName(),
			c.Points,
		))},
		{render.Centred(RenderClan(c.Clan))},
	}
	for i, l := range c.CalcLines(stealing) {
		row := []interface{}{render.Markup(fmt.Sprintf(
			"%d.",
			i+1,
		), render.Gray, false)}
		row = append(row, l.RenderRow()...)
		cells = append(cells, []interface{}{
			render.Table([][]interface{}{row}, 0, 2),
		})
	}
	return cells
}

func (l Line) RenderRow() []interface{} {
	row := []interface{}{}
	for _, s := range l.Symbols {
		row = append(row, RenderDie(s))
	}
	if l.Infantry > 0 {
		row = append(row, RenderInf(l.Infantry))
	}
	return row
}

func (l Line) String() string {
	return render.Table([][]interface{}{l.RenderRow()}, 0, 2)
}
