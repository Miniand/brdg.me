package starship_catan

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	playerNum, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	opponentNum := (playerNum + 1) % 2
	// Current turn
	cells := [][]interface{}{
		[]interface{}{
			Bold("Current turn:"),
			g.RenderName(playerNum),
		},
	}
	switch g.Phase {
	case PhaseChooseSector:
		if len(g.PlayerBoards[playerNum].LastSectors) > 0 {
			cells = append(
				cells,
				[]interface{}{
					Bold("Last sectors"),
					strings.Join(Itoas(g.PlayerBoards[playerNum].LastSectors), " "),
				},
			)
		}
	case PhaseFlight:
		if g.FlightCards.Len() > 0 {
			card, _ := g.FlightCards.Pop()
			cells = append(
				cells,
				[]interface{}{
					Bold("Current planet:"),
					fmt.Sprintf("%s", card),
				},
				[]interface{}{
					Bold("Current sector:"),
					strconv.Itoa(g.CurrentSector),
				},
				[]interface{}{
					Bold("Moves left:"),
					strconv.Itoa(g.RemainingMoves()),
				},
				[]interface{}{
					Bold("Actions left:"),
					strconv.Itoa(g.RemainingActions()),
				},
			)
		}
	case PhaseTradeAndBuild:
		cells = append(
			cells,
			[]interface{}{
				Bold("Post trades remaining:"),
				strconv.Itoa(g.RemainingTrades()),
			},
			[]interface{}{
				Bold("Player trades remaining:"),
				strconv.Itoa(g.RemainingPlayerTrades()),
			},
		)
	}
	t := render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Resources
	cells = [][]interface{}{
		[]interface{}{
			Bold("Resource"), Bold(g.RenderName(playerNum)), Bold(g.RenderName(opponentNum)),
			" ", // Column spacing
			Bold("Resource"), Bold(g.RenderName(playerNum)), Bold(g.RenderName(opponentNum)),
		},
		g.ResourceTableDoubleRow(ResourceFood, ResourceColonyShip, playerNum),
		g.ResourceTableDoubleRow(ResourceFuel, ResourceTradeShip, playerNum),
		g.ResourceTableDoubleRow(ResourceCarbon, ResourceBooster, playerNum),
		g.ResourceTableDoubleRow(ResourceOre, ResourceCannon, playerNum),
		g.ResourceTableRow(ResourceTrade, playerNum),
		DoubleRow(
			g.ResourceTableRow(ResourceScience, playerNum),
			[]interface{}{
				fmt.Sprintf(`{{c "red"}}{{b}}medals{{_b}}{{_c}}`),
				Bold(strconv.Itoa(g.PlayerBoards[playerNum].Medals())),
				strconv.Itoa(g.PlayerBoards[opponentNum].Medals()),
			},
		),
		DoubleRow(
			[]interface{}{"", "", ""},
			[]interface{}{
				fmt.Sprintf(`{{c "green"}}{{b}}diplomacy{{_b}}{{_c}}`),
				Bold(strconv.Itoa(g.PlayerBoards[playerNum].DiplomatPoints())),
				strconv.Itoa(g.PlayerBoards[opponentNum].DiplomatPoints()),
			},
		),
		DoubleRow(
			g.ResourceTableRow(ResourceAstro, playerNum),
			[]interface{}{
				fmt.Sprintf(`{{c "blue"}}{{b}}VP{{_b}}{{_c}}`),
				Bold(strconv.Itoa(g.PlayerBoards[playerNum].VictoryPoints())),
				strconv.Itoa(g.PlayerBoards[opponentNum].VictoryPoints()),
			},
		),
	}
	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Adventure cards
	buf.WriteString("{{b}}Adventure cards{{_b}}\n")
	cells = [][]interface{}{
		[]interface{}{Bold("#"), Bold("Planet"), Bold("Description")},
	}
	for i, c := range g.CurrentAdventureCards() {
		ac := c.(Adventurer)
		cells = append(cells, []interface{}{
			strconv.Itoa(i + 1),
			AdventurePlanetString(ac.Planet()),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, ac.Text()),
		})
	}
	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Modules
	cells = [][]interface{}{{
		"{{b}}Module{{_b}}",
		g.RenderName(playerNum),
		g.RenderName(opponentNum),
		"{{b}}Description{{_b}}",
	}}
	for _, m := range Modules {
		cells = append(cells, []interface{}{
			ModuleNames[m],
			RenderModuleLevel(g.PlayerBoards[playerNum].Modules[m]),
			RenderModuleLevel(g.PlayerBoards[opponentNum].Modules[m]),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, ModuleSummaries[m]),
		})
	}
	t = render.Table(cells, 0, 2)
	buf.WriteString(t)
	buf.WriteString(fmt.Sprintf(
		`
{{b}}Upgrade cost: L1{{_b}} (%s), {{b}}L2{{_b}} (%s)`,
		ModuleTransaction(1).LoseString(),
		ModuleTransaction(2).LoseString(),
	))
	// Cards
	for _, p := range []int{playerNum, opponentNum} {
		buf.WriteString("\n\n")
		strs := []string{
			fmt.Sprintf("%s {{b}}cards{{_b}}", g.RenderName(p)),
		}
		for _, c := range g.PlayerBoards[p].Colonies {
			strs = append(strs, fmt.Sprintf("%s", c))
		}
		for _, c := range g.PlayerBoards[p].TradingPosts {
			strs = append(strs, fmt.Sprintf("%s", c))
		}
		buf.WriteString(strings.Join(strs, "\n"))
	}
	return buf.String(), nil
}

func RenderModuleLevel(level int) string {
	switch level {
	case 0:
		return fmt.Sprintf(`{{c "gray"}}0{{_c}}`)
	case 2:
		return fmt.Sprintf(`{{b}}2{{_b}}`)
	default:
		return strconv.Itoa(level)
	}
}

func (g *Game) ResourceTableRow(resource, player int) []interface{} {
	opponent := (player + 1) % 2
	return []interface{}{
		RenderResource(resource),
		Bold(strconv.Itoa(g.PlayerBoards[player].Resources[resource])),
		fmt.Sprintf(`{{c "gray"}}%d{{_c}}`,
			g.PlayerBoards[opponent].Resources[resource]),
	}
}

func DoubleRow(row1, row2 []interface{}) []interface{} {
	row1 = append(row1, "")
	row1 = append(row1, row2...)
	return row1
}

func (g *Game) ResourceTableDoubleRow(resource1, resource2, player int) []interface{} {
	return DoubleRow(
		g.ResourceTableRow(resource1, player),
		g.ResourceTableRow(resource2, player),
	)
}

func RenderMoney(amount int) string {
	return fmt.Sprintf(`{{b}}{{c "green"}}$%d{{_c}}{{_b}}`, amount)
}

func RenderResource(resource int) string {
	if _, ok := ResourceColours[resource]; !ok {
		log.Fatalf(
			"There is no resource colour for %s (%d)",
			ResourceNames[resource],
			resource,
		)
	}
	return fmt.Sprintf(
		`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`,
		ResourceColours[resource],
		ResourceNames[resource],
	)
}

func RenderResourceAmount(resource, amount int) string {
	switch resource {
	case ResourceAstro:
		return RenderMoney(amount)
	default:
		return fmt.Sprintf("%d %s", amount, RenderResource(resource))
	}
}

func RenderResources(resources []int) string {
	strs := make([]string, len(resources))
	for i, r := range resources {
		strs[i] = RenderResource(r)
	}
	return strings.Join(strs, ", ")
}

func Bold(s string) string {
	return fmt.Sprintf("{{b}}%s{{_b}}", s)
}
