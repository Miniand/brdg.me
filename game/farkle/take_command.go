package farkle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/die"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type TakeCommand struct{}

func (tc TakeCommand) Name() string { return "take" }

func (tc TakeCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("can't find player")
	}
	if !g.CanTake(pNum) {
		return "", errors.New("can't take at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify what to take")
	}
	takeString := args[0]
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
			"%s is not valid, please check the usage examples",
			takeString,
		))
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
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took %s for {{b}}%d{{_b}} points",
		render.PlayerName(g.Player, g.Players[g.Player]),
		RenderDice(take),
		score,
	)))
	return "", nil
}

func (tc TakeCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	scoreStrings := []string{}
	for _, s := range AvailableScores(g.RemainingDice) {
		scoreStrings = append(scoreStrings, fmt.Sprintf("     take %s",
			s.Description()))
	}
	return fmt.Sprintf(
		"{{b}}take #{{_b}} to take some dice for points.  You can:\n%s",
		strings.Join(scoreStrings, "\n"))
}

func (g *Game) CanTake(player int) bool {
	return player == g.Player && !g.IsFinished() &&
		len(AvailableScores(g.RemainingDice)) > 0
}
