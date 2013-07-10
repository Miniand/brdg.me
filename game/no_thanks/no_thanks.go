package no_thanks

import (
	"errors"
	"math/rand"
	"time"
)

type Game struct {
	Players         []string
	PlayerHands     map[string][]int
	PlayerChips     map[string]int
	RemainingCards  []int
	CurrentlyMoving string
}

func (g *Game) PlayerAction(player, action string, params []string) error {
	return nil
}

func (g *Game) Name() string {
	return "No Thanks"
}

func (g *Game) Identifier() string {
	return "no_thanks"
}

func (g *Game) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (g *Game) Decode([]byte) error {
	return nil
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 || len(players) > 5 {
		return errors.New("No Thanks requires between 2 and 5 players")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Players = players
	g.InitCards()
	g.InitPlayerChips()
	g.CurrentlyMoving = g.Players[r.Int()%len(g.Players)]
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
	return []string{}
}

func (g *Game) AllCards() []int {
	cards := make([]int, 33)
	for i := 3; i <= 35; i++ {
		cards[i-3] = i
	}
	return cards
}

func (g *Game) InitCards() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	cardPool := g.AllCards()
	picked := map[int]bool{}
	g.RemainingCards = make([]int, 24)
	for i := 0; i < 24; i++ {
		c := cardPool[r.Int()%24]
		for picked[c] {
			c = cardPool[r.Int()%24]
		}
		picked[c] = true
		g.RemainingCards[i] = c
	}
}

func (g *Game) InitPlayerChips() {
	g.PlayerChips = map[string]int{}
	for _, p := range g.Players {
		g.PlayerChips[p] = 11
	}
}
