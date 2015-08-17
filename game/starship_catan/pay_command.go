package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type PayCommand struct{}

func (c PayCommand) Name() string { return "pay" }

func (c PayCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.PayRansom(p)
}

func (c PayCommand) Usage(player string, context interface{}) string {
	return "{{b}}pay{{_b}} to pay the ransom"
}

func (g *Game) CanPayRansom(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight ||
		g.FlightCards.Len() == 0 || g.LosingModule {
		return false
	}
	card, _ := g.FlightCards.Pop()
	pirateCard, ok := card.(PirateCard)
	return ok && pirateCard.Ransom <= g.PlayerBoards[player].Resources[ResourceAstro]
}

func (g *Game) PayRansom(player int) error {
	if !g.CanPayRansom(player) {
		return errors.New("you aren't able to pay the ransom")
	}
	card, _ := g.FlightCards.Pop()
	pirateCard, ok := card.(PirateCard)
	if !ok {
		return errors.New("card isn't a pirate card")
	}
	g.PlayerBoards[player].Resources[ResourceAstro] -= pirateCard.Ransom
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s paid a ransom of %s`,
		g.RenderName(player),
		RenderMoney(pirateCard.Ransom),
	)))
	g.CardFinished = true
	return nil
}
