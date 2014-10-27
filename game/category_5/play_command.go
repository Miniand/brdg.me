package category_5

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("play", 1, input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanPlay(pNum)
}

func (c PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which card to play")
	}
	cardNum, err := strconv.Atoi(a[0])
	if err != nil {
		return "", err
	}

	return "", g.Play(pNum, Card(cardNum))
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play #{{_b}} to play a card, eg. {{b}}play 47{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return !g.Resolving && g.Plays[player] == 0
}

func (g *Game) Play(player int, card Card) error {
	if !g.CanPlay(player) {
		return errors.New("you can't play at the moment")
	}

	var ok bool
	g.Hands[player], ok = RemoveCard(g.Hands[player], card)
	if !ok {
		return errors.New("you don't have that card")
	}

	g.Plays[player] = card

	// Check if everyone had played
	for p, _ := range g.Players {
		if g.Plays[p] == 0 {
			// Some people haven't played yet
			return nil
		}
	}
	g.ResolvePlays()
	return nil
}
