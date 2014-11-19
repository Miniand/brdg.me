package king_of_tokyo

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type LeaveCommand struct{}

func (c LeaveCommand) Parse(input string) []string {
	return command.ParseNamedCommand("leave", input)
}

func (c LeaveCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanLeave(pNum)
}

func (c LeaveCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Leave(pNum)
}

func (c LeaveCommand) Usage(player string, context interface{}) string {
	return "{{b}}leave{{_b}} to leave Tokyo"
}

func (g *Game) CanLeave(player int) bool {
	return g.CanStay(player)
}

func (g *Game) Leave(player int) error {
	if !g.CanLeave(player) {
		return errors.New("you can't call leave at the moment")
	}
	g.Tokyo[g.CurrentFleeingLoc] = TokyoEmpty
	g.PostStayOrLeave()
	return nil
}
