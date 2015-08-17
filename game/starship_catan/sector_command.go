package starship_catan

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type SectorCommand struct{}

func (c SectorCommand) Name() string { return "sector" }

func (c SectorCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a sector")
	}
	s, err := strconv.Atoi(args[0])
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
	switch g.PlayerBoards[g.CurrentPlayer].Modules[ModuleSensor] {
	case 1:
		g.Peeking, g.SectorCards[g.CurrentSector] =
			g.SectorCards[g.CurrentSector].PopN(2)
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s is using the sensor module to peek at 2 cards`,
			g.RenderName(player),
		)))
	case 2:
		g.Peeking, g.SectorCards[g.CurrentSector] =
			g.SectorCards[g.CurrentSector].PopN(3)
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s is using the sensor module to peek at 3 cards`,
			g.RenderName(player),
		)))
	default:
		return g.NextSectorCard()
	}
	return nil
}
