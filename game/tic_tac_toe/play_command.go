package tic_tac_toe

import (
	"errors"
	"regexp"
	"strings"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return regexp.MustCompile(`(?im)^\s*([a-i])\s*$`).FindStringSubmatch(input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return !g.IsFinished() && player == g.CurrentlyMoving
}

// Make an action for the specified player
func (c PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	action := strings.ToLower(args[1])
	if g.CurrentlyMoving != player {
		return "", errors.New("Not your turn")
	}
	if !regexp.MustCompile("^[abcdefghi]$").MatchString(action) {
		return "", errors.New("Your action must be a letter between a - i")
	}
	var x, y int
	switch action {
	case "a":
		x = 0
		y = 0
	case "b":
		x = 1
		y = 0
	case "c":
		x = 2
		y = 0
	case "d":
		x = 0
		y = 1
	case "e":
		x = 1
		y = 1
	case "f":
		x = 2
		y = 1
	case "g":
		x = 0
		y = 2
	case "h":
		x = 1
		y = 2
	case "i":
		x = 2
		y = 2
	}
	err := g.MarkCellForPlayer(player, x, y)
	if err != nil {
		return "", err
	}
	g.NextPlayer()
	return "", nil
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "Send a letter between {{b}}a - i{{_b}} to play in that square"
}
