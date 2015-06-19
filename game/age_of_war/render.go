package age_of_war

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	layout := [][]interface{}{
		{render.Centred(render.Bold("Current roll"))},
		{render.Centred(strings.Join(RenderDice(g.CurrentRoll), "   "))},
		{},
	}

	if g.CurrentlyAttacking != -1 {
		layout = append(layout, [][]interface{}{
			{render.Centred(render.Bold("Currently attacking"))},
			{},
			{render.Centred(g.RenderCastle(g.CurrentlyAttacking))},
			{},
		}...)
	}

	layout = append(layout, [][]interface{}{
		{},
		{render.Centred(render.Bold("Castles"))},
		{},
	}...)
	layout = append(layout, []interface{}{
		render.Centred(g.RenderCastles()),
	})

	scores := g.Scores()
	scoreStrs := make([]string, len(g.Players))
	for p := range g.Players {
		scoreStrs[p] = fmt.Sprintf(
			"%s: {{b}}%d{{_b}}",
			g.PlayerName(p),
			scores[p],
		)
	}
	layout = append(layout, [][]interface{}{
		{},
		{render.Centred(render.Bold("Scores"))},
		{render.Centred(strings.Join(scoreStrs, "   "))},
	}...)

	return render.Table(layout, 0, 0), nil
}

func (g *Game) RenderCastles() string {
	cells := [][]interface{}{}
	row := []interface{}{}
	lastClan := -1
	conqueredClans := map[int]bool{}
	for i, c := range Castles {
		if lastClan != -1 && c.Clan != lastClan && len(row) > 0 {
			cells = append(cells, []interface{}{render.Centred(render.Table(
				[][]interface{}{row}, 0, 6))})
			row = []interface{}{}
		}
		conquered, ok := conqueredClans[c.Clan]
		if !ok {
			var conqueredBy int
			conquered, conqueredBy = g.ClanConquered(c.Clan)
			conqueredClans[c.Clan] = conquered
			if conquered {
				cells = append(cells, []interface{}{
					render.Centred(fmt.Sprintf(
						"%s has been conquered by %s for {{b}}%d{{_b}} points",
						RenderClan(c.Clan),
						g.PlayerName(conqueredBy),
						ClanSetPoints[c.Clan],
					)),
				})
			}
		}
		if !conquered {
			row = append(row, g.RenderCastle(i))
		}
		lastClan = c.Clan
	}
	if len(row) > 0 {
		cells = append(cells, []interface{}{render.Centred(render.Table(
			[][]interface{}{row}, 0, 6))})
	}
	return render.Table(cells, 1, 6)
}

func (g *Game) RenderCastle(cIndex int) string {
	c := Castles[cIndex]
	cells := [][]interface{}{
		{render.Centred(fmt.Sprintf(
			"%s (%d)",
			c.RenderName(),
			c.Points,
		))},
	}
	if g.Conquered[cIndex] {
		cells = append(cells, []interface{}{render.Centred(fmt.Sprintf(
			"(%s)",
			g.PlayerName(g.CastleOwners[cIndex]),
		))})
	}
	for i, l := range c.CalcLines(g.Conquered[cIndex]) {
		row := []interface{}{render.Markup(fmt.Sprintf(
			"%d.",
			i+1,
		), render.Gray, false)}
		if g.CurrentlyAttacking == cIndex && g.CompletedLines[i] {
			row = append(row, render.Colour("complete", render.Gray))
		} else {
			row = append(row, l.RenderRow()...)
		}
		cells = append(cells, []interface{}{
			render.Table([][]interface{}{row}, 0, 2),
		})
	}
	return render.Table(cells, 0, 0)
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
