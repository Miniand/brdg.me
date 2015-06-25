package love_letter

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players                []string
	Log                    *log.Log
	Round                  int
	Deck, Removed, Discard []int
	Hands                  [][]int
	Points                 []int
	CurrentPlayer          int
	Eliminated             []bool
	Protected              []bool
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
	}
}

func (g *Game) Name() string {
	return "Love Letter"
}

func (g *Game) Identifier() string {
	return "love_letter"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	l := len(players)
	if l < 2 || l > 4 {
		return errors.New("only for 2 to 4 players")
	}
	g.Players = players
	g.Log = log.New()
	g.Points = make([]int, l)
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.Round++
	g.Eliminated = make([]bool, len(g.Players))
	g.Protected = make([]bool, len(g.Players))
	deck := helper.IntShuffle(Deck)
	remove := 1
	l := len(g.Players)
	if l == 2 {
		remove = 4
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Starting round %d, {{b}}removing %d %s{{_b}}",
		g.Round,
		remove,
		helper.Plural(remove, "card"),
	)))
	g.Deck, g.Removed = deck[remove:], deck[:remove]
	g.Hands = make([][]int, l)
	for p := range g.Players {
		g.Hands[p] = []int{}
		g.DrawCard(p)
	}
	g.StartTurn()
}

func (g *Game) StartTurn() {
	g.Protected[g.CurrentPlayer] = false
	if len(g.Deck) == 0 {
		g.EndRound()
	} else {
		g.DrawCard(g.CurrentPlayer)
	}
}

func (g *Game) NextPlayer() {
	for {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		if !g.Eliminated[g.CurrentPlayer] {
			break
		}
	}
	g.StartTurn()
}

func (g *Game) DiscardCardLog(player, card int) {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded %s",
		g.RenderName(player),
		RenderCard(card),
	)))
	g.DiscardCard(player, card)
}

func (g *Game) DiscardCard(player, card int) {
	g.Hands[player] = helper.IntRemove(card, g.Hands[player], 1)
	g.Discard = append(g.Discard, card)
	if card == Princess {
		g.Eliminate(player)
	}
}

func (g *Game) Eliminate(player int) {
	if g.Eliminated[player] {
		return
	}

	g.Eliminated[player] = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s has been eliminated from this round",
		g.RenderName(player),
	)))
	for len(g.Hands[player]) > 0 {
		g.DiscardCardLog(player, g.Hands[player][0])
	}

	numRemaining := 0
	for p := range g.Players {
		if !g.Eliminated[p] {
			numRemaining++
		}
	}
	if numRemaining <= 1 {
		g.EndRound()
	}
}

func (g *Game) EndRound() {
	// Scoring
	g.StartRound()
}

func (g *Game) DrawCard(player int) {
	var card int
	if len(g.Deck) > 0 {
		card, g.Deck = g.Deck[0], g.Deck[1:]
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s drew a card from the draw pile, {{b}}%d{{_b}} remaining",
			g.RenderName(player),
			len(g.Deck),
		)))
	} else {
		card = g.Removed[0]
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s drew a card from the removed cards",
			g.RenderName(player),
		)))
	}
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"You drew %s",
		RenderCard(card),
	), []string{g.Players[player]}))
	g.Hands[player] = append(g.Hands[player], card)
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
	return helper.StringInStrings(player, g.Players)
}

func (g *Game) AvailableTargets(forPlayer int) []int {
	targets := []int{}
	for p := range g.Players {
		if p != forPlayer && !g.Eliminated[p] && !g.Protected[p] {
			targets = append(targets, p)
		}
	}
	return targets
}
