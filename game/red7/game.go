package red7

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players []string
	Log     *log.Log

	Finished bool

	CurrentPlayer int
	HasPlayed     bool

	Deck        []int
	DiscardPile []int
	Hands       [][]int
	Palettes    [][]int
	ScoredCards [][]int
	Eliminated  []bool
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
		DiscardCommand{},
		DoneCommand{},
	}
}

func (g *Game) Name() string {
	return "Red7"
}

func (g *Game) Identifier() string {
	return "red7"
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
		return errors.New("only 2-4 players")
	}
	g.Players = players
	g.Log = log.New()

	g.Deck = append([]int{}, Deck...)
	g.DiscardPile = []int{}
	g.Hands = make([][]int, l)
	g.Palettes = make([][]int, l)
	g.ScoredCards = make([][]int, l)

	g.StartRound()

	return nil
}

func (g *Game) StartRound() {
	l := len(g.Players)

	// Add hands and palettes back to the deck.
	for p := range g.Players {
		g.Deck = append(g.Deck, g.Hands[p]...)
		g.Deck = append(g.Deck, g.Palettes[p]...)
	}
	g.Deck = append(g.Deck, g.DiscardPile...)
	g.DiscardPile = []int{}
	g.Hands = make([][]int, l)
	g.Palettes = make([][]int, l)
	g.Eliminated = make([]bool, l)

	if len(g.Deck) < l*8 {
		// End of the game, not enough cards to deal new hand.
		g.EndGame()
		return
	}

	g.Deck = helper.IntShuffle(g.Deck)

	// Deal hands and new palettes.
	for p := range g.Players {
		g.Draw(p, 7)
		g.Palettes[p] = g.Deck[0:1]
		g.Deck = g.Deck[1:]
	}

	// Starting player is to the left of the current leader.
	leader, _ := g.Leader()
	g.CurrentPlayer = g.NextPlayer(leader)
	g.StartTurn()
}

func (g *Game) Draw(player, n int) {
	l := len(g.Deck)
	if l == 0 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s didn't draw from the deck as there are no cards left",
			g.PlayerName(player),
		)))
		return
	}
	if l < n {
		n = l
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s drew {{b}}%d{{_b}} cards from the deck",
		g.PlayerName(player),
		n,
	)))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"You drew %s",
		strings.Join(RenderCards(g.Deck[:n]), " "),
	), []string{g.Players[player]}))
	g.Hands[player] = append(g.Hands[player], g.Deck[:n]...)
	g.Deck = g.Deck[n:]
}

func (g *Game) StartTurn() {
	g.HasPlayed = false
	if len(g.Hands[g.CurrentPlayer]) == 0 {
		g.Eliminate(g.CurrentPlayer, " for not having any cards left")
		g.EndTurn()
	}
}

func (g *Game) EndTurn() {
	if !g.HasPlayed {
		g.Eliminate(g.CurrentPlayer, " for not playing or discarding")
	} else if leader, _ := g.Leader(); leader != g.CurrentPlayer {
		g.Eliminate(g.CurrentPlayer, " for not being the leader at the end of their turn")
	}

	if len(g.RemainingPlayers()) == 1 {
		g.EndRound()
		return
	}

	g.CurrentPlayer = g.NextPlayer(g.CurrentPlayer)
	g.StartTurn()
}

func (g *Game) EndRound() {
	leader, leaderPalette := g.Leader()
	g.ScoredCards[leader] = append(g.ScoredCards[leader], leaderPalette...)
	g.Palettes[leader], _ = helper.IntSliceSub(g.Palettes[leader], leaderPalette)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s won the round with %s for {{b}}%d{{_b}} points, now on {{b}}%d{{_b}} points",
		g.PlayerName(leader),
		strings.Join(RenderCards(helper.IntSort(leaderPalette)), " "),
		Points(leaderPalette),
		g.PlayerPoints(leader),
	)))

	endPoints := EndPoints(len(g.Players))
	for p := range g.Players {
		if g.PlayerPoints(p) >= endPoints {
			g.EndGame()
			return
		}
	}

	g.StartRound()
}

func (g *Game) EndGame() {
	g.Log.Add(log.NewPublicMessage(render.Bold("It is the end of the game")))
	g.Finished = true
}

func (g *Game) Eliminate(player int, message string) {
	g.Eliminated[player] = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s has been eliminated%s",
		g.PlayerName(player),
		message,
	)))
}

func (g *Game) Leader() (leader int, palette []int) {
	return g.LeaderWithSuit(g.CurrentRule())
}

func (g *Game) LeaderWithSuit(suit int) (leader int, palette []int) {
	playerMap := map[int]int{}
	palettes := [][]int{}
	for p := range g.Players {
		if g.Eliminated[p] {
			continue
		}
		playerMap[len(palettes)] = p
		palettes = append(palettes, SuitRules[suit](g.Palettes[p]))
	}
	lIndex, palette := Leader(palettes)
	leader = playerMap[lIndex]
	return
}

func (g *Game) CurrentRule() int {
	l := len(g.DiscardPile)
	if l == 0 {
		return SuitRed
	}
	return g.DiscardPile[l-1] & SuitMask
}

func (g *Game) NextPlayer(from int) int {
	l := len(g.Players)
	n := (from + 1) % l
	for {
		if n == from || !g.Eliminated[n] {
			break
		}
		n = (n + 1) % l
	}
	return n
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Finished
}

func (g *Game) Winners() []string {
	if !g.Finished {
		return []string{}
	}
	winnerPoints := 0
	winners := []string{}
	for p := range g.Players {
		points := g.PlayerPoints(p)
		if points > winnerPoints {
			winnerPoints = points
			winners = []string{}
		}
		if points == winnerPoints {
			winners = append(winners, g.Players[p])
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for p, pName := range g.Players {
		if pName == player {
			return p, true
		}
	}
	return 0, false
}

func (g *Game) PlayerPoints(player int) int {
	return Points(g.ScoredCards[player])
}

func (g *Game) RemainingPlayers() []int {
	rem := []int{}
	for p := range g.Players {
		if !g.Eliminated[p] {
			rem = append(rem, p)
		}
	}
	return rem
}

func EndPoints(players int) int {
	return 50 - players*5
}
