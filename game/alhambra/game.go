package alhambra

import (
	"encoding/gob"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	PhaseAction = iota
	PhasePlace
	PhaseEnd
)

type Game struct {
	Players []string
	Log     *log.Log

	CurrentPlayer int
	Phase         int
	Round         int

	Boards []PlayerBoard

	Cards, CardPile, DiscardPile card.Deck
	Tiles, TileBag               []Tile
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		BuyCommand{},
	}
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

	g.Round = 1

	g.CardPile = Deck().Shuffle()
	g.Cards = g.DrawCards(4)
	g.DiscardPile = card.Deck{}

	g.TileBag = ShuffleTiles(Tiles)
	g.Tiles, g.TileBag = g.TileBag[:4], g.TileBag[4:]

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
			c = g.DrawCards(1)[0]
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

func (g *Game) DrawCards(n int) card.Deck {
	var c card.Card
	if n <= 0 {
		return card.Deck{}
	}
	if g.CardPile.Len() == 0 {
		if g.DiscardPile.Len() == 0 {
			return card.Deck{}
		}
		g.CardPile = g.DiscardPile.Shuffle()
		g.DiscardPile = card.Deck{}
	}
	c, g.CardPile = g.CardPile.Pop()
	switch c.(type) {
	case Card:
		return (card.Deck{c}).PushMany(g.DrawCards(n - 1))
	case ScoringCard:
		g.ScoreRound()
		return g.DrawCards(n)
	default:
		panic("Unknown card")
	}
}

func (g *Game) ScoreRound() {
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

var ErrCouldNotFindPlayer = errors.New("could not find player")

func (g *Game) PlayerNum(player string) (int, bool) {
	for pNum, p := range g.Players {
		if p == player {
			return pNum, true
		}
	}
	return 0, false
}

func (g *Game) NextPhase() {
	switch g.Phase {
	case PhaseAction:
		g.PlacePhase()
	case PhasePlace:
		g.NextPlayer()
	}
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	g.ActionPhase()
}

func (g *Game) ActionPhase() {
	g.Phase = PhaseAction
}

func (g *Game) PlacePhase() {
	g.Phase = PhasePlace
	if len(g.Boards[g.CurrentPlayer].Place) == 0 {
		g.NextPhase()
	}
}

var ParseCardRegexp = regexp.MustCompile(`(?i)^([a-z])([0-9]+)$`)

func ParseCard(input string) (Card, error) {
	matches := ParseCardRegexp.FindStringSubmatch(input)
	if matches == nil {
		return Card{}, errors.New("cards must be a letter followed by a number, such as R10")
	}
	currency, err := helper.MatchStringInStringMap(matches[1], CurrencyAbbr)
	if err != nil {
		return Card{}, err
	}
	value, err := strconv.Atoi(matches[2])
	if err != nil {
		return Card{}, err
	}
	return Card{currency, value}, nil
}
