package farkle

import (
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/die"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"strings"
)

type TakeCommand struct{}

func (tc TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("take", 1, input)
}

func (tc TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return player == g.Players[g.Player] && !g.IsFinished() &&
		len(AvailableScores(g.RemainingDice)) > 0
}

func (tc TakeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	a := command.ExtractNamedCommandArgs(args)
	takeString := a[0]
	g := context.(*Game)
	if player != g.Players[g.Player] {
		return "", errors.New("It's not your turn")
	}
	if g.IsFinished() {
		return "", errors.New("The game is already finished")
	}
	take, err := die.ValueStringToDice(takeString)
	if err != nil {
		return "", err
	}
	// Check that it's a valid value string and get the points
	score := 0
	for _, s := range Scores() {
		if die.DiceEquals(take, s.Dice) {
			score = s.Value
			break
		}
	}
	if score == 0 {
		return "", errors.New(fmt.Sprintf(
			"%s is not valid, see below for valid dice to take:\n%s",
			takeString, strings.Join(ScoreStrings(), "\n")))
	}
	// Check that we've actually got the dice
	isIn, remaining := die.DiceInDice(take, g.RemainingDice)
	if !isIn {
		return "", errors.New("You don't have the correct dice for " +
			takeString)
	}
	g.TurnScore += score
	g.TakenThisRoll = true
	g.RemainingDice = remaining
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took %s for %d points",
		render.PlayerName(g.Player, g.Players[g.Player]),
		RenderDice(take), score)))
	return "", nil
}

func (tc TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take #{{_b}} to take some dice for points, eg. {{b}}take 222{{_b}} or {{b}}take 5{{_b}}"
}
