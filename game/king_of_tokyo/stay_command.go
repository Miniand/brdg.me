package king_of_tokyo

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type StayCommand struct{}

func (c StayCommand) Parse(input string) []string {
	return command.ParseNamedCommand("stay", input)
}

func (c StayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanStay(pNum)
}

func (c StayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Stay(pNum)
}

func (c StayCommand) Usage(player string, context interface{}) string {
	return "{{b}}stay{{_b}} to stay in Tokyo"
}

func (g *Game) CanStay(player int) bool {
	return g.Phase == PhaseFlee && g.Tokyo[g.CurrentFleeingLoc] == player
}

func (g *Game) Stay(player int) error {
	if !g.CanStay(player) {
		return errors.New("you can't call stay at the moment")
	}
	g.PostStayOrLeave()
	return nil
}

func (g *Game) PostStayOrLeave() {
	p := g.Tokyo[g.CurrentFleeingLoc]
	g.Boards[p].Health -= g.AttackDamage
	if g.Boards[p].Health < 0 {
		g.Boards[p].Health = 0
		// Remove from Tokyo
		if pLoc := g.PlayerLocation(p); pLoc != LocationOutside {
			g.Tokyo[pLoc] = TokyoEmpty
		}
	}
	g.NextFleeingLoc()
}
