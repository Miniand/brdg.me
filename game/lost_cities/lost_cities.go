package lost_cities

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	Players         []string
	CurrentlyMoving string
	StartPlayer     string
	Board           [5][10]int
	PlayerHands     map[string][]int
	DrawStack       []int
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Lost Cities requires 2 spieler")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(r)
	g.Players = players
	//g.InitCards()
	//g.CurrentlyMoving = g.Players[r.Int()%len(g.Players)]
	return nil
}

func (g *Game) PlayerAction(player, action string, params []string) error {
	return nil
}

func (g *Game) Name() string {
	return "Lost Cities"
}

func (g *Game) Identifier() string {
	return "lost_cities"
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
