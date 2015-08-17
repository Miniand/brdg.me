package starship_catan

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type CompleteCommand struct{}

func (c CompleteCommand) Name() string { return "complete" }

func (c CompleteCommand) Call(
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
		return "", errors.New("you must specify an adventure number")
	}
	adventure, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("you must specify an adventure number")
	}

	return "", g.Complete(p, adventure)
}

func (c CompleteCommand) Usage(player string, context interface{}) string {
	return "{{b}}complete #{{_b}} to complete an adventure, eg. {{b}}complete 2{{_b}}"
}

func (g *Game) CanComplete(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight {
		return false
	}
	c, _ := g.FlightCards.Pop()
	_, ok := c.(AdventurePlanetCard)
	return ok && g.CurrentAdventureCards().Len() > 0
}

func (g *Game) Complete(player, adventure int) error {
	if !g.CanComplete(player) {
		return errors.New("you can't complete an adventure at the moment")
	}
	if adventure <= 0 {
		return errors.New("the adventure number must be above 0")
	}
	current := g.CurrentAdventureCards()
	if l := len(current); adventure > l {
		return fmt.Errorf("the adventure number can't be higher than %d", l)
	}

	c, _ := g.FlightCards.Pop()
	apc := c.(AdventurePlanetCard)
	ac := current[adventure-1].(Adventurer)
	if ac.Planet() != apc.Name {
		return errors.New("it is not the correct planet to complete that card")
	}

	if err := ac.Complete(player, g); err != nil {
		return err
	}

	g.MarkCardActioned()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s completed a mission on %s - {{c "gray"}}%s{{_c}}`,
		g.RenderName(player),
		apc,
		ac.Text(),
	)))
	g.RemoveAdventureCard = adventure
	if g.GainResources == nil {
		g.Completed()
	}
	return nil
}
