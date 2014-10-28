package category_5

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type ChooseCommand struct{}

func (c ChooseCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("choose", 1, input)
}

func (c ChooseCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanChoose(pNum)
}

func (c ChooseCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which row to choose")
	}
	row, err := strconv.Atoi(a[0])
	if err != nil {
		return "", err
	}

	return "", g.Choose(pNum, row)
}

func (c ChooseCommand) Usage(player string, context interface{}) string {
	return "{{b}}choose #{{_b}} to choose which row to take, eg. {{b}}choose 4{{_b}}"
}

func (g *Game) CanChoose(player int) bool {
	return g.Resolving && g.ChoosePlayer == player
}

func (g *Game) Choose(player, row int) error {
	if !g.CanChoose(player) {
		return errors.New("you can't choose at the moment")
	}

	if row < 1 || row > 4 {
		return errors.New("the row must be between 1 and 4")
	}
	row -= 1

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s and chose to take row {{b}}%d{{_b}} for {{b}}%d points{{_b}}",
		g.RenderName(player),
		g.Plays[player],
		row+1,
		CardsHeads(g.Board[row]),
	)))

	g.PlayerCards[player] = append(g.PlayerCards[player], g.Board[row]...)
	g.Board[row] = []Card{g.Plays[player]}
	g.Plays[player] = 0

	g.ResolvePlays()
	return nil
}
