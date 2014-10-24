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
	cells := [][]string{
		[]string{
			Bold("Current turn:"),
			g.RenderName(playerNum),
		},
	}
	switch g.Phase {
	case PhaseChooseSector:
		if len(g.PlayerBoards[playerNum].LastSectors) > 0 {
			cells = append(
				cells,
				[]string{
					Bold("Last sectors"),
					strings.Join(Itoas(g.PlayerBoards[playerNum].LastSectors), " "),
				},
			)
		}
	case PhaseFlight:
		card, _ := g.FlightCards.Pop()
		cells = append(
			cells,
			[]string{
				Bold("Current planet:"),
				fmt.Sprintf("%s", card),
			},
			[]string{
				Bold("Current sector:"),
				strconv.Itoa(g.CurrentSector),
			},
			[]string{
				Bold("Moves left:"),
				strconv.Itoa(g.RemainingMoves()),
			},
			[]string{
				Bold("Actions left:"),
				strconv.Itoa(g.RemainingActions()),
			},
		)
	case PhaseTradeAndBuild:
		cells = append(
			cells,
			[]string{
				Bold("Post trades remaining:"),
				strconv.Itoa(g.RemainingTrades()),
			},
			[]string{
				Bold("Player trades remaining:"),
				strconv.Itoa(g.RemainingPlayerTrades()),
			},
		)
	}
	t, err := render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Resources
	cells = [][]string{
		[]string{Bold("Resource"), Bold(g.RenderName(playerNum)), Bold(g.RenderName(opponentNum))},
		g.ResourceTableRow(ResourceFood, playerNum),
		g.ResourceTableRow(ResourceFuel, playerNum),
		g.ResourceTableRow(ResourceCarbon, playerNum),
		g.ResourceTableRow(ResourceOre, playerNum),
		g.ResourceTableRow(ResourceTrade, playerNum),
		g.ResourceTableRow(ResourceScience, playerNum),
		[]string{},
		g.ResourceTableRow(ResourceAstro, playerNum),
		[]string{},
		g.ResourceTableRow(ResourceColonyShip, playerNum),
		g.ResourceTableRow(ResourceTradeShip, playerNum),
		[]string{},
		g.ResourceTableRow(ResourceBooster, playerNum),
		g.ResourceTableRow(ResourceCannon, playerNum),
	}
	t, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Adventure cards
	buf.WriteString("{{b}}Adventure cards{{_b}}\n")
	cells = [][]string{
		[]string{Bold("#"), Bold("Planet"), Bold("Description")},
	}
	for i, c := range g.CurrentAdventureCards() {
		ac := c.(Adventurer)
		cells = append(cells, []string{
			strconv.Itoa(i + 1),
			AdventurePlanetString(ac.Planet()),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, ac.Text()),
		})
	}
	t, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(t)
	buf.WriteString("\n\n")
	// Modules
	cells = [][]string{{
		"{{b}}Module{{_b}}",
		g.RenderName(playerNum),
		g.RenderName(opponentNum),
		"{{b}}Description{{_b}}",
	}}
	for _, m := range Modules {
		cells = append(cells, []string{
			ModuleNames[m],
			RenderModuleLevel(g.PlayerBoards[playerNum].Modules[m]),
			RenderModuleLevel(g.PlayerBoards[opponentNum].Modules[m]),
			fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, ModuleSummaries[m]),
		})
	}
	t, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(t)
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

func (g *Game) ResourceTableRow(resource, player int) []string {
	opponent := (player + 1) % 2
	return []string{
		RenderResource(resource),
		Bold(strconv.Itoa(g.PlayerBoards[player].Resources[resource])),
		fmt.Sprintf(`{{c "gray"}}%d{{_c}}`,
			g.PlayerBoards[opponent].Resources[resource]),
	}
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
