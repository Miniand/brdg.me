package scommand

import (
	"bytes"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

type ListCommand struct{}

func (c ListCommand) Name() string { return "list" }

func (c ListCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	output := bytes.NewBufferString("Available games:")
	for gName, g := range game.RawCollection() {
		output.WriteString("\n  ")
		output.WriteString(gName)
		output.WriteString(" (")
		output.WriteString(g.Name())
		output.WriteString(")")
	}
	return output.String(), nil
}

func (c ListCommand) Usage(player string, context interface{}) string {
	return "{{b}}list{{_b}} to get a list of available games"
}

func CanList(player string, gm *model.GameModel) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || ok && u.Unsubscribed {
		return false
	}
	return gm == nil || gm.IsFinished
}
