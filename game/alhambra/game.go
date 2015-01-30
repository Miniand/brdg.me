package alhambra

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players []string
	Log     *log.Log

	CurrentPlayer int

	Boards []PlayerBoard

	Cards card.Deck
	Tiles []Tile

	DrawCards card.Deck
	DrawTiles []Tile
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Alhambra"
}

func (g *Game) Identifier() string {
	return "alhambra"
}

func RegisterGobTypes() {
	gob.Register(Card{})
}

func (g *Game) Encode() ([]byte, error) {
	RegisterGobTypes()
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	RegisterGobTypes()
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if l := len(players); l < 2 || l > 6 {
		return errors.New("Alhambra requires between 2 and 6 players")
	}
	g.Players = players
	g.Log = log.New()

	g.DrawCards = Deck().Shuffle()
	g.DrawTiles = ShuffleTiles(Tiles)

	g.Boards = make([]PlayerBoard, len(g.Players))
	var (
		c                             card.Card
		minPlayer, minCards, minValue int
	)
	for pNum, _ := range g.Players {
		g.Boards[pNum] = NewPlayerBoard()
		cards := 0
		value := 0
		cardStrs := []string{}
		for value < 20 {
			c, g.DrawCards = g.DrawCards.Pop()
			cards++
			value += c.(Card).Value
			cardStrs = append(cardStrs, c.(Card).String())
			g.Boards[pNum].Cards = g.Boards[pNum].Cards.Push(c)
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s drew %s",
			g.PlayerName(pNum),
			render.CommaList(cardStrs),
		)))
		if pNum == 0 || cards < minCards || cards == minCards && value < minValue {
			minPlayer = pNum
			minCards = cards
			minValue = value
		}
	}

	g.CurrentPlayer = minPlayer
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s is the starting player as they got the fewest cards",
		g.PlayerName(minPlayer),
	)))

	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for pNum, p := range g.Players {
		if p == player {
			return pNum, true
		}
	}
	return 0, false
}
