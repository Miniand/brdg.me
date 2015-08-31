package age_of_war

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type LineCommand struct{}

func (c LineCommand) Name() string { return "line" }

func (c LineCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a line to complete")
	}

	line, err := strconv.Atoi(args[0])
	if err != nil || line <= 0 {
		return "", errors.New("the line must be a number greater than 0")
	}
	return "", g.Line(pNum, line-1)
}

func (c LineCommand) Usage(player string, context interface{}) string {
	return "{{b}}line #{{_b}} to complete a line in the castle you are attacking, eg. {{b}}line 2{{_b}}"
}

func (g *Game) CanLine(player int) bool {
	return g.CurrentPlayer == player && g.CurrentlyAttacking != -1
}

func (g *Game) Line(player, line int) error {
	if !g.CanLine(player) {
		return errors.New("unable to complete a line right now")
	}
	lines := Castles[g.CurrentlyAttacking].CalcLines(
		g.Conquered[g.CurrentlyAttacking],
	)
	if line < 0 || line >= len(lines) {
		return errors.New("that is not a valid line")
	}
	if g.CompletedLines[line] {
		return errors.New("that line has already been completed")
	}
	canAfford, with := lines[line].CanAfford(g.CurrentRoll)
	if !canAfford {
		return errors.New("cannot afford that line")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s completed %s with {{b}}%d{{_b}} %s",
		g.PlayerName(player),
		lines[line].String(),
		with,
		helper.Plural(with, "die"),
	)))
	g.CompletedLines[line] = true
	// Check end of turn first in case they completed the castle.
	if !g.CheckEndOfTurn() {
		g.Roll(len(g.CurrentRoll) - with)
		g.CheckEndOfTurn()
	}
	return nil
}
