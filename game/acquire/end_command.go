package acquire

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type EndCommand struct{}

func (c EndCommand) Name() string { return "end" }

func (c EndCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanEnd(pNum) {
		return "", errors.New("you can't end at the moment")
	}
	g.FinalTurn = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} triggered the end of the game at the end of their turn`,
		g.RenderPlayer(g.CurrentPlayer))))
	if !g.PlayerCanAffordShares(g.CurrentPlayer) {
		// Player can't buy anything so just end the game.
		g.NextPlayer()
	}
	return "", nil
}

func (c EndCommand) Usage(player string, context interface{}) string {
	return `{{b}}end{{_b}} to end the game after your current turn`
}
