package splendor

import (
	"bytes"
	"fmt"
	"strconv"
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
	Diamond:  "Diam",
	Sapphire: "Saph",
	Emerald:  "Emer",
	Ruby:     "Ruby",
	Onyx:     "Onyx",
	Gold:     "Gold",
	Prestige: "VP",
}

func GemStrings() map[int]string {
	strs := map[int]string{}
	for _, g := range Gems {
		strs[g] = ResourceStrings[g]
	}
	return strs
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	pb := g.PlayerBoards[pNum]
	bonuses := pb.Bonuses()

	output := bytes.NewBuffer([]byte{})

	// Nobles
	nobleHeader := []interface{}{""}
	nobleRow := []interface{}{render.Colour(fmt.Sprintf(
		"Nobles (%s each)",
		render.Bold(RenderResourceColour(3, Prestige)),
	), render.Gray)}
	for i, n := range g.Nobles {
		nobleHeader = append(nobleHeader,
			render.Centred(render.Colour(i+1, render.Gray)))
		nobleRow = append(nobleRow, RenderAmount(n.Cost))
	}
	table := [][]interface{}{
		nobleHeader,
		nobleRow,
	}
	output.WriteString(render.Table(table, 0, 2))
	output.WriteString("\n\n")

	// Board
	longestRow := 0
	for _, r := range g.Board {
		if l := len(r); l > longestRow {
			longestRow = l
		}
	}
	if l := len(g.PlayerBoards[pNum].Reserve); l > longestRow {
		longestRow = l
	}
	header := []interface{}{""}
	for i := 0; i < longestRow; i++ {
		header = append(header, render.Centred(
			render.Colour(fmt.Sprintf("{{b}}%c{{_b}}", 'A'+i), render.Gray)))
	}
	table = [][]interface{}{header}
	for l, r := range g.Board {
		upper := []interface{}{
			render.Colour(fmt.Sprintf("Level {{b}}%d{{_b}}", l+1), render.Gray),
		}
		lower := []interface{}{
			"",
		}
		for _, c := range r {
			upperBuf := bytes.NewBuffer([]byte{})
			if bonuses.CanAfford(c.Cost) {
				upperBuf.WriteString(render.Markup("✔ ", render.Green, true))
			} else if g.PlayerBoards[pNum].CanAfford(c.Cost) {
				upperBuf.WriteString(render.Markup("✔ ", render.Yellow, true))
			}
			upperBuf.WriteString(RenderCardBonusVP(c))
			upper = append(upper, render.Centred(upperBuf.String()))
			lower = append(lower, render.Centred(RenderAmount(c.Cost)))
		}
		table = append(table, upper, lower, []interface{}{})
	}
	upper := []interface{}{
		render.Colour("Level {{b}}4{{_b}}", render.Gray),
	}
	lower := []interface{}{
		render.Colour("Reserved", render.Gray),
	}
	for _, c := range g.PlayerBoards[pNum].Reserve {
		upperBuf := bytes.NewBuffer([]byte{})
		if g.PlayerBoards[pNum].CanAfford(c.Cost) {
			upperBuf.WriteString(render.Markup("✔ ", render.Green, true))
		}
		upperBuf.WriteString(RenderCardBonusVP(c))
		upper = append(upper, render.Centred(upperBuf.String()))
		lower = append(lower, render.Centred(RenderAmount(c.Cost)))
	}
	table = append(table, upper, lower)
	output.WriteString(render.Table(table, 0, 3))
	output.WriteString("\n\n\n")

	// Tokens
	tableHeader := []interface{}{""}
	yourTokenRow := []interface{}{render.Bold("You have")}
	availTokenRow := []interface{}{render.Bold("Tokens left")}
	for _, gem := range append(Gems, Gold) {
		var yourTokenCell string
		tableHeader = append(tableHeader,
			render.Centred(render.Bold(RenderResourceColour(ResourceAbbr[gem], gem))))
		yourTokenCell = render.Bold(strconv.Itoa(bonuses[gem]))
		if gem == Gold {
			yourTokenCell = render.Bold(strconv.Itoa(pb.Tokens[Gold]))
		} else if pb.Tokens[gem] > 0 {
			yourTokenCell = fmt.Sprintf(
				"{{b}}%d{{_b}} (%d+%d)",
				bonuses[gem]+pb.Tokens[gem],
				bonuses[gem],
				pb.Tokens[gem],
			)
		}
		yourTokenRow = append(yourTokenRow, render.Centred(yourTokenCell))
		availTokenRow = append(availTokenRow,
			render.Centred(strconv.Itoa(g.Tokens[gem])))
	}
	table = [][]interface{}{
		tableHeader,
		yourTokenRow,
		availTokenRow,
	}
	output.WriteString(render.Table(table, 0, 3))
	output.WriteString("\n\n\n")

	// Player table
	header = []interface{}{""}
	for _, gem := range Gems {
		header = append(header, render.Centred(render.Bold(
			RenderResourceColour(ResourceAbbr[gem], gem))))
	}
	header = append(
		header,
		render.Centred(render.Bold(
			RenderResourceColour(ResourceAbbr[Gold], Gold))),
		render.Centred(render.Bold(
			render.Colour("Res", render.Cyan))),
		render.Centred(render.Bold(
			RenderResourceColour(ResourceAbbr[Prestige], Prestige))),
	)
	table = [][]interface{}{header}
	for p, _ := range g.Players {
		bold := p == pNum
		pb := g.PlayerBoards[p]
		bonuses := pb.Bonuses()
		row := []interface{}{g.RenderName(p)}
		for _, gem := range Gems {
			gemBuf := bytes.NewBufferString(strconv.Itoa(bonuses[gem]))
			if n := pb.Tokens[gem]; n > 0 {
				gemBuf.WriteString(fmt.Sprintf("+%d", n))
			}
			row = append(row, render.Centred(render.Markup(
				gemBuf.String(), "", bold)))
		}
		row = append(
			row,
			render.Centred(render.Markup(pb.Tokens[Gold], "", bold)),
			render.Centred(render.Markup(len(pb.Reserve), "", bold)),
			render.Centred(render.Markup(pb.Prestige(), "", bold)),
		)
		table = append(table, row)
	}
	output.WriteString(render.Table(table, 0, 2))

	return output.String(), nil
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
	return render.Bold(strings.Join(parts, " "))
}

func RenderAmount(a Amount) string {
	parts := []string{}
	for _, r := range Resources {
		if a[r] > 0 {
			parts = append(parts, render.Bold(RenderResourceColour(a[r], r)))
		}
	}
	return strings.Join(parts, render.Colour("-", render.Gray))
}

func RenderNobleHeader(n Noble) string {
	return RenderResourceColour(n.Prestige, Prestige)
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func RenderCard(c Card) string {
	return render.Bold(fmt.Sprintf(
		"%s (%s)",
		RenderCardBonusVP(c),
		RenderAmount(c.Cost),
	))
}

func RenderNoble(n Noble) string {
	return render.Bold(fmt.Sprintf(
		"%s (%s)",
		RenderNobleHeader(n),
		RenderAmount(n.Cost),
	))
}
