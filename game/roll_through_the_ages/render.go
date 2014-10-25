package roll_through_the_ages

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func (g *Game) RenderForPlayer(player string) (string, error) {
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
	t, err := render.Table([][]string{diceRow, numberRow}, 0, 2)
	if err != nil {
		return "", err
	}
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
					`{{b}}{{c "%s"}}X{{_c}}{{_b}} `,
					render.PlayerColour(p),
				), g.Boards[p].CityProgress+1),
				strings.Repeat(
					`{{c "gray"}}.{{_c}} `,
					MaxCityProgress-g.Boards[p].CityProgress,
				),
			),
		})
	}
	t, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
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
		for p, pn := range g.Players {
			cell := `{{c "gray"}}.{{_c}}`
			if g.Boards[p].Developments[d] {
				cell = fmt.Sprintf(
					`{{b}}{{c "%s"}}X{{_c}}{{_b}}`,
					render.PlayerColour(p),
				)
			}
			row = append(row, fmt.Sprintf(
				`%s%s`,
				strings.Repeat(" ", len(pn)/2+1),
				cell,
			))
		}
		row = append(row, []string{
			fmt.Sprintf(" %d", dv.Cost),
			fmt.Sprintf(" %d", dv.Points),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, dv.Effect),
		}...)
		cells = append(cells, row)
	}
	t, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
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
		for p, pn := range g.Players {
			var cell string
			switch {
			case g.Boards[p].Monuments[m] == 0:
				cell = `{{c "gray"}}.{{_c}}`
			case g.Boards[p].MonumentBuiltFirst[m]:
				cell = fmt.Sprintf(
					`{{b}}{{c "%s"}}X{{_c}}{{_b}}`,
					render.PlayerColour(p),
				)
			case g.Boards[p].Monuments[m] == mv.Size:
				cell = fmt.Sprintf(
					`{{c "%s"}}x{{_c}}`,
					render.PlayerColour(p),
				)
			default:
				cell = fmt.Sprintf(
					`{{c "%s"}}%d{{_c}}`,
					render.PlayerColour(p),
					g.Boards[p].Monuments[m],
				)
			}
			row = append(row, fmt.Sprintf(
				`%s%s`,
				strings.Repeat(" ", len(pn)/2+1),
				cell,
			))
		}
		row = append(row, []string{
			fmt.Sprintf(" %d", mv.Size),
			fmt.Sprintf("{{b}}%d{{_b}}/%d", mv.Points, mv.SubsequentPoints()),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, mv.Effect),
		}...)
		cells = append(cells, row)
	}
	t, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(t)
	return buf.String(), nil
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

func Bold(s string) string {
	return fmt.Sprintf("{{b}}%s{{_b}}", s)
}
