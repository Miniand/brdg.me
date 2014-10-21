package starship_catan

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

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
	if g.Phase == PhaseFlight {
		card, _ := g.FlightCards.Pop()
		cells = append(
			cells,
			[]string{
				Bold("Current planet:"),
				fmt.Sprintf("%v", card),
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
	return buf.String(), nil
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

func Bold(s string) string {
	return fmt.Sprintf("{{b}}%s{{_b}}", s)
}
