package roll_through_the_ages

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	// Calculate name widths now as we use them quite a lot
	nameWidths := map[int]int{}
	for p, _ := range g.Players {
		nameWidths[p] = render.StrLen(g.RenderName(p))
	}
	buf := bytes.NewBuffer([]byte{})
	// Dice
	diceRow := []string{}
	numberRow := []string{}
	for i, d := range g.RolledDice {
		diceString := DiceStrings[d]
		diceRow = append(diceRow, Bold(RenderDice(d)))
		numberRow = append(numberRow, fmt.Sprintf(
			`%s{{c "gray"}}%d{{_c}}`,
			strings.Repeat(" ", utf8.RuneCountInString(diceString)/2),
			i+1,
		))
	}
	for _, d := range g.KeptDice {
		diceRow = append(diceRow, RenderDice(d))
	}
	buf.WriteString("{{b}}Dice{{_b}} {{c \"gray\"}}(F: food, W: worker, G: good, C: coin, X: skull){{_c}}\n")
	t := render.Table([][]string{diceRow, numberRow}, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Cities
	buf.WriteString("{{b}}Cities{{_b}} {{c \"gray\"}}(number of dice and food used per turn){{_c}}\n")
	cityHeaderBuf := bytes.NewBufferString(fmt.Sprintf(
		"{{b}}%d{{_b}}", BaseCitySize))
	last := 0
	for i, n := range CityLevels {
		cityHeaderBuf.WriteString(fmt.Sprintf(
			`%s{{b}}%d{{_b}}`,
			strings.Repeat(" ", (n-last-1)*2+1),
			BaseCitySize+i+1,
		))
		last = n
	}
	cells := [][]string{{"{{b}}Player{{_b}}", cityHeaderBuf.String()}}
	for p, _ := range g.Players {
		cells = append(cells, []string{
			g.RenderName(p),
			fmt.Sprintf(
				"%s%s",
				strings.Repeat(fmt.Sprintf(
					`%s `,
					RenderX(p, p == pNum),
				), g.Boards[p].CityProgress+1),
				strings.Repeat(
					`{{c "gray"}}.{{_c}} `,
					MaxCityProgress-g.Boards[p].CityProgress,
				),
			),
		})
	}
	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Developments
	header := []string{Bold("Development")}
	for p, _ := range g.Players {
		header = append(header, g.RenderName(p))
	}
	header = append(header, []string{
		Bold("Cost"),
		Bold("Pts"),
		Bold("Effect"),
	}...)
	cells = [][]string{header}
	for _, d := range Developments {
		dv := DevelopmentValues[d]
		row := []string{strings.Title(dv.Name)}
		for p, _ := range g.Players {
			cell := `{{c "gray"}}.{{_c}}`
			if g.Boards[p].Developments[d] {
				cell = RenderX(p, pNum == p)
			}
			row = append(row, render.Centre(cell, nameWidths[p]))
		}
		row = append(row, []string{
			fmt.Sprintf(" %d", dv.Cost),
			fmt.Sprintf(" %d", dv.Points),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, dv.Effect),
		}...)
		cells = append(cells, row)
	}
	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Monuments
	header = []string{Bold("Monument")}
	for p, _ := range g.Players {
		header = append(header, g.RenderName(p))
	}
	header = append(header, []string{
		Bold("Size"),
		Bold("Pts"),
		Bold("Effect"),
	}...)
	cells = [][]string{header}
	for _, m := range g.Monuments() {
		mv := MonumentValues[m]
		row := []string{strings.Title(mv.Name)}
		for p, _ := range g.Players {
			var cell string
			switch {
			case g.Boards[p].Monuments[m] == 0:
				cell = `{{c "gray"}}.{{_c}}`
			case g.Boards[p].Monuments[m] == mv.Size:
				cell = RenderX(p, g.Boards[p].MonumentBuiltFirst[m])
			default:
				cell = fmt.Sprintf(
					`{{c "%s"}}%d{{_c}}`,
					render.PlayerColour(p),
					g.Boards[p].Monuments[m],
				)
			}
			row = append(row, render.Centre(cell, nameWidths[p]))
		}
		row = append(row, []string{
			fmt.Sprintf(" %d", mv.Size),
			fmt.Sprintf("{{b}}%d{{_b}}/%d", mv.Points, mv.SubsequentPoints()),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, mv.Effect),
		}...)
		cells = append(cells, row)
	}
	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Resources
	header = []string{Bold("Resource")}
	for p, _ := range g.Players {
		header = append(header, g.RenderName(p))
	}
	cells = [][]string{header}
	for _, good := range GoodsReversed() {
		row := []string{RenderGoodName(good)}
		for p, _ := range g.Players {
			num := g.Boards[p].Goods[good]
			cell := Colour(".", "gray")
			if num > 0 {
				cell = Markup(
					fmt.Sprintf("%d (%d)", num, GoodValue(good, num)),
					render.PlayerColour(p),
					p == pNum,
				)
			}
			row = append(row, render.Centre(cell, nameWidths[p]))
		}
		cells = append(cells, row)
	}
	row := []string{Bold("total")}
	for p, _ := range g.Players {
		cell := Markup(
			fmt.Sprintf("%d (%d)", g.Boards[p].GoodsNum(), g.Boards[p].GoodsValue()),
			render.PlayerColour(p),
			p == pNum,
		)
		row = append(row, render.Centre(cell, nameWidths[p]))
	}
	cells = append(cells, row, []string{})

	row = []string{FoodName}
	for p, _ := range g.Players {
		cell := Markup(
			strconv.Itoa(g.Boards[p].Food),
			render.PlayerColour(p),
			p == pNum,
		)
		row = append(row, render.Centre(cell, nameWidths[p]))
	}
	cells = append(cells, row)
	for p, _ := range g.Players {
		cell := Markup(
			strconv.Itoa(g.Boards[p].Disasters),
			render.PlayerColour(p),
			p == pNum,
		)
		row = append(row, render.Centre(cell, nameWidths[p]))
	}
	row = []string{DisasterName}
	cells = append(cells, row)
	row = []string{Bold("score")}
	for p, _ := range g.Players {
		cell := Markup(
			strconv.Itoa(g.Boards[p].Score()),
			render.PlayerColour(p),
			p == pNum,
		)
		row = append(row, render.Centre(cell, nameWidths[p]))
	}
	cells = append(cells, row)

	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	return buf.String(), nil
}

func RenderX(player int, strong bool) string {
	x := "x"
	if strong {
		x = "X"
	}
	return Markup(x, render.PlayerColour(player), strong)
}

func RenderDice(dice int) string {
	diceString := DiceStrings[dice]
	for v, col := range DiceValueColours {
		diceString = strings.Replace(diceString, v, fmt.Sprintf(
			`{{c "%s"}}%s{{_c}}`,
			col,
			v,
		), -1)
	}
	return diceString
}

func RenderGoodName(good int) string {
	return fmt.Sprintf(
		`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`,
		GoodColours[good],
		GoodStrings[good],
	)
}

var FoodName = `{{b}}{{c "green"}}food{{_c}}{{_b}}`
var DisasterName = `{{b}}{{c "red"}}disaster{{_c}}{{_b}}`

func Markup(s string, colour string, bold bool) string {
	if colour != "" {
		s = Colour(s, colour)
	}
	if bold {
		s = Bold(s)
	}
	return s
}

func Colour(s, colour string) string {
	return fmt.Sprintf(`{{c "%s"}}%s{{_c}}`, colour, s)
}

func Bold(s string) string {
	return fmt.Sprintf("{{b}}%s{{_b}}", s)
}
