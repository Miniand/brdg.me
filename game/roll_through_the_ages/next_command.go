package roll_through_the_ages

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type NextCommand struct{}

func (c NextCommand) Name() string { return "next" }

func (c NextCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Next(pNum)
}

func (c NextCommand) Usage(player string, context interface{}) string {
	return "{{b}}next{{_b}} to continue to the next phase of your turn"
}

func (g *Game) CanNext(player int) bool {
	return player == g.CurrentPlayer && ContainsInt(g.Phase, []int{
		PhasePreserve,
		PhaseRoll,
		PhaseExtraRoll,
		PhaseInvade,
		PhaseBuild,
		PhaseTrade,
		PhaseBuy,
	})
}

func (g *Game) Next(player int) error {
	if !g.CanNext(player) {
		return errors.New("you can't next at the moment")
	}
	g.NextPhase()
	return nil
}
