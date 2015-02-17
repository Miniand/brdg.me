package zombie_dice

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

var DiceFaceStrings = map[int]string{
	Brain:      "Brain",
	Shotgun:    "Shot",
	Footprints: "Run",
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	output := bytes.NewBuffer([]byte{})
	cupStr := render.Markup("None", render.Gray, false)
	if len(g.Cup) > 0 {
		counts := make([]int, 3)
		for _, d := range g.Cup {
			counts[ColourOrder[d.Colour]]++
		}
		parts := []string{}
		for _, c := range Colours {
			if counts[ColourOrder[c]] > 0 {
				parts = append(parts, render.Markup(
					fmt.Sprintf("%d %s", counts[ColourOrder[c]], c),
					c,
					true,
				))
			}
		}
		cupStr = render.CommaList(parts)
	}
	output.WriteString(render.Table([][]interface{}{
		{
			render.RightAligned("Brains:"),
			render.Bold(strconv.Itoa(g.RoundBrains)),
		},
		{
			render.RightAligned("Shots:"),
			render.Bold(strconv.Itoa(g.RoundShotguns)),
		},
		{
			render.RightAligned("Runners:"),
			g.CurrentRoll.String(),
		},
		{
			render.RightAligned("Kept:"),
			g.Kept.String(),
		},
		{
			render.RightAligned("In cup:"),
			cupStr,
		},
	}, 0, 2))
	output.WriteString("\n\n\n{{b}}Scores:{{_b}}\n")
	cells := [][]interface{}{}
	for p, _ := range g.Players {
		cells = append(cells, []interface{}{
			render.RightAligned(g.PlayerName(p)),
			render.Bold(strconv.Itoa(g.Scores[p])),
		})
	}
	output.WriteString(render.Table(cells, 0, 2))
	output.WriteString(`


{{b}}Notes:{{_b}}
{{c "green"}}{{b}}Green{{_b}}{{_c}} dice have {{b}}3 Brain{{_b}}, {{b}}2 Run{{_b}} and {{b}}1 Shot{{_b}}.
{{c "yellow"}}{{b}}Yellow{{_b}}{{_c}} dice have {{b}}2 Brain{{_b}}, {{b}}2 Run{{_b}} and {{b}}2 Shot{{_b}}.
{{c "red"}}{{b}}Red{{_b}}{{_c}} dice have {{b}}1 Brain{{_b}}, {{b}}2 Run{{_b}} and {{b}}3 Shot{{_b}}.`)
	return output.String(), nil
}

func (d DiceResult) String() string {
	return render.Markup(DiceFaceStrings[d.Face], d.Colour, true)
}

func (drl DiceResultList) String() string {
	parts := make([]string, len(drl))
	for i, dr := range drl {
		parts[i] = dr.String()
	}
	return strings.Join(parts, " ")
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
