package starship_catan

import (
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
	opponentNum := (playerNum + 1) % 2
	cells := [][]string{
		[]string{Bold("Resource"), Bold(g.RenderName(playerNum)), Bold(g.RenderName(opponentNum))},
		g.ResourceTableRow(ResourceAstro, playerNum),
		g.ResourceTableRow(ResourceColonyShip, playerNum),
		g.ResourceTableRow(ResourceTradeShip, playerNum),
		g.ResourceTableRow(ResourceBooster, playerNum),
		g.ResourceTableRow(ResourceCannon, playerNum),
		[]string{},
		g.ResourceTableRow(ResourceFood, playerNum),
		g.ResourceTableRow(ResourceFuel, playerNum),
		g.ResourceTableRow(ResourceCarbon, playerNum),
		g.ResourceTableRow(ResourceOre, playerNum),
		g.ResourceTableRow(ResourceTrade, playerNum),
		g.ResourceTableRow(ResourceScience, playerNum),
	}
	t, err := render.Table(cells, 0, 2)
	return t, err
}

func (g *Game) ResourceTableRow(resource, player int) []string {
	opponent := (player + 1) % 2
	return []string{
		RenderResource(resource),
		Bold(strconv.Itoa(g.PlayerBoards[player].Resources[resource])),
		strconv.Itoa(g.PlayerBoards[opponent].Resources[resource]),
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
