package starship_catan

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type GainCommand struct{}

func (c GainCommand) Name() string { return "gain" }

func (c GainCommand) Call(
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
	if err != nil || len(args) == 0 {
		return "", errors.New("you must specify what resource to gain")
	}
	r, err := ParseResource(args[0])
	if err != nil {
		return "", err
	}
	found := false
	for _, gr := range g.GainResources {
		if gr == r {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf(
			`You aren't able to gain {{b}}%s{{_b}} at the moment`,
			ResourceNames[r])
	}
	g.GainResource(p, r)
	return "", g.Gained(p)
}

func (c GainCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	resources := make([]string, len(g.GainResources))
	for i, r := range g.GainResources {
		resources[i] = fmt.Sprintf("{{b}}%s{{_b}}", ResourceNames[r])
	}
	return fmt.Sprintf(
		"{{b}}gain ##{{_b}} to gain a resource.  Enter as much of the resource name as needed to uniquely identify it.  Eg. {{b}}gain sci{{_b}}\nYou can gain: %s",
		strings.Join(resources, ", "))
}

func (g *Game) CanGain(player int) bool {
	return g.GainPlayer == player && g.GainResources != nil
}

func (g *Game) GainOne(player int, resources []int) {
	if g.GainResources != nil {
		// Add it to the queue
		g.GainQueue = append(g.GainQueue, resources)
		return
	}
	if len(resources) == 0 {
		g.Gained(player)
	}
	canProduce := g.PlayerBoards[player].FitTransaction(
		TransactionFromResources(resources)).Resources()
	switch len(canProduce) {
	case 0:
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s did not gain a resource, all full", g.RenderName(player))))
		g.Gained(player)
	case 1:
		g.GainResource(player, canProduce[0])
		g.Gained(player)
	default:
		g.GainPlayer = player
		g.GainResources = canProduce
	}
}

func (g *Game) Gained(player int) error {
	g.GainResources = nil
	if len(g.GainQueue) > 0 {
		// There's still some gains in the queue, kick off the next one.
		resources := g.GainQueue[0]
		g.GainQueue = g.GainQueue[1:]
		g.GainOne(player, resources)
		return nil
	}
	switch g.Phase {
	case PhaseProduce:
		if player == g.CurrentPlayer {
			g.Produce((g.CurrentPlayer + 1) % 2)
		} else {
			g.Phase = PhaseChooseSector
		}
	case PhaseFlight:
		c, _ := g.FlightCards.Pop()
		switch c.(type) {
		case AdventurePlanetCard:
			g.Completed()
		}
	}
	return nil
}
