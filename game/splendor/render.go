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
	Diamond:  "Di",
	Sapphire: "Sa",
	Emerald:  "Em",
	Ruby:     "Ru",
	Onyx:     "On",
	Gold:     "Go",
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

	output := bytes.NewBuffer([]byte{})

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
				render.Markup(RenderAmount(c.Cost), "", canAfford)))
		}
		table = append(table, upper, lower, []interface{}{})
	}
	upper := []interface{}{
		render.Colour("Level 4", render.Gray),
	}
	lower := []interface{}{
		render.Colour("Reserved", render.Gray),
	}
	for _, c := range g.PlayerBoards[pNum].Reserve {
		canAfford := g.PlayerBoards[pNum].CanAfford(c.Cost)
		upper = append(upper, render.Centred(
			render.Markup(RenderCardBonusVP(c), "", canAfford)))
		lower = append(lower, render.Centred(
			render.Markup(RenderAmount(c.Cost), "", canAfford)))
	}
	table = append(table, upper, lower)
	output.WriteString(render.Table(table, 0, 3))
	output.WriteString("\n\n")

	// Nobles
	nobleHeader := []interface{}{""}
	nobleRow := []interface{}{render.Colour(fmt.Sprintf(
		"Nobles (%s each)",
		render.Bold(RenderResourceColour(3, Prestige)),
	), render.Gray)}
	for i, n := range g.Nobles {
		nobleHeader = append(nobleHeader,
			render.Centred(render.Colour(i+1, render.Gray)))
		nobleRow = append(nobleRow, render.Bold(RenderAmount(n.Cost)))
	}
	table = [][]interface{}{
		nobleHeader,
		nobleRow,
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
	tokensRow := []interface{}{render.Markup("Tokens", render.Gray, true)}
	for _, gem := range Gems {
		tokensRow = append(tokensRow, render.Centred(strconv.Itoa(g.Tokens[gem])))
	}
	tokensRow = append(
		tokensRow,
		render.Centred(strconv.Itoa(g.Tokens[Gold])),
		render.Centred(render.Colour("-", render.Gray)),
		render.Centred(render.Colour("-", render.Gray)),
	)
	table = [][]interface{}{header, tokensRow}
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
	return strings.Join(parts, " ")
}

func RenderAmount(a Amount) string {
	parts := []string{}
	for _, r := range Resources {
		if a[r] > 0 {
			parts = append(parts, RenderResourceColour(a[r], r))
		}
	}
	return strings.Join(parts, "")
}

func RenderNobleHeader(n Noble) string {
	return RenderResourceColour(n.Prestige, Prestige)
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
