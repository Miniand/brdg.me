package starship_catan

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type SectorCommand struct{}

func (c SectorCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("sector", 1, input)
}

func (c SectorCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanSector(p)
}

func (c SectorCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	s, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("sector must be a number between 1 and 4")
	}
	g.Sector(p, s)
	return "", nil
}

func (c SectorCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	p, _ := g.ParsePlayer(player)
	lastSectorMsg := ""
	if len(g.PlayerBoards[p].LastSectors) != 0 {
		lastSectorMsg = fmt.Sprintf(
			`Your last sector was {{b}}sector %d{{_b}}.  `,
			g.PlayerBoards[p].LastSectors[0])
	}
	return fmt.Sprintf(
		"{{b}}sector #{{_b}} to choose which sector to travel through, between 1 and 4.  %sEg. {{b}}sector 3{{_b}}", lastSectorMsg)
}

func (g *Game) CanSector(player int) bool {
	return g.Phase == PhaseChooseSector && g.CurrentPlayer == player
}

func (g *Game) Sector(player, sector int) error {
	if !g.CanSector(player) {
		return errors.New("you can't choose a sectore at the moment")
	}
	if sector < 1 || sector > 4 {
		return errors.New("sector must be between 1 and 4")
	}
	g.Phase = PhaseFlight
	g.CurrentSector = sector
	g.FlightActions = map[int]bool{}
	return g.NextSectorCard()
}
