package starship_catan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	PutOnTop = iota
	PutOnBottom
)

var PutStrings = map[int]string{
	PutOnTop:    "top",
	PutOnBottom: "bottom",
}

type PutCommand struct{}

func (c PutCommand) Name() string { return "put" }

func (c PutCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 2 {
		return "", errors.New("you must specify which card and where to put it")
	}

	num, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("the first argument must be a positive number")
	}

	on, err := helper.MatchStringInStringMap(args[1], PutStrings)
	if err != nil {
		return "", err
	}
	return "", g.Put(p, num, on)
}

func (c PutCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	cards := make([]string, len(g.Peeking))
	for i, c := range g.Peeking {
		str := ""
		switch t := c.(type) {
		case FullStringer:
			str = t.FullString()
		case fmt.Stringer:
			str = t.String()
		}
		cards[i] = fmt.Sprintf("%d. %s", i+1, str)
	}
	return fmt.Sprintf(
		`{{b}}put # top/bottom{{_b}} to put a card on the top or the bottom of the pile, eg. {{b}}put 1 bottom{{_b}}
The cards are:
%s`,
		strings.Join(cards, "\n"))
}

func (g *Game) CanPut(player int) bool {
	return g.CurrentPlayer == player && g.Peeking.Len() > 0
}

func (g *Game) Put(player, num, on int) error {
	if !g.CanPut(player) {
		return errors.New("you can't put cards at the moment")
	}
	if num < 1 || num > g.Peeking.Len() {
		return errors.New("you must specify the number of one of the listed cards")
	}
	if on != PutOnTop && on != PutOnBottom {
		return errors.New("invalid on value")
	}

	index := num - 1
	c := g.Peeking[index]
	g.Peeking = append(g.Peeking[:index], g.Peeking[index+1:]...)
	switch on {
	case PutOnTop:
		g.SectorCards[g.CurrentSector] =
			g.SectorCards[g.CurrentSector].Push(c)
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s put a card on the top of the pile`,
			g.RenderName(player),
		)))
	case PutOnBottom:
		g.SectorCards[g.CurrentSector] =
			g.SectorCards[g.CurrentSector].Unshift(c)
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s put a card on the bottom of the pile`,
			g.RenderName(player),
		)))
	}
	if g.Peeking.Len() == 0 {
		return g.NextSectorCard()
	}
	return nil
}
