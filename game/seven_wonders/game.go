package seven_wonders

import (
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

	Round   int
	Hands   []card.Deck
	Discard card.Deck
	Actions []Actioner

	Cards         []card.Deck
	Coins         []int
	VictoryTokens []int
	DefeatTokens  []int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		BuildCommand{},
		DealCommand{},
		DiscardCommand{},
	}
}

func (g *Game) Name() string {
	return "7 Wonders"
}

func (g *Game) Identifier() string {
	return "7_wonders"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	pLen := len(players)
	if pLen < 3 || pLen > 7 {
		return errors.New("7 Wonders is 3 to 7 player")
	}
	g.Players = players
	g.Log = log.New()

	g.Discard = card.Deck{}

	g.Cards = make([]card.Deck, pLen)
	g.Coins = make([]int, pLen)
	g.VictoryTokens = make([]int, pLen)
	g.DefeatTokens = make([]int, pLen)
	for i := 0; i < pLen; i++ {
		g.Cards[i] = card.Deck{}
		g.Coins[i] = 3
	}

	g.StartRound(1)

	return nil
}

func (g *Game) StartRound(round int) {
	if round < 4 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"It is now {{b}}round %d{{_b}}",
			round,
		)))
	}
	g.Round = round
	players := len(g.Players)
	switch round {
	case 1:
		g.DealHands(DeckAge1(players).Shuffle())
	case 2:
		g.DealHands(DeckAge2(players).Shuffle())
	case 3:
		g.DealHands(DeckAge3(players).Shuffle())
	case 4:
		// End of game
	}
}

func (g *Game) DealHands(cards card.Deck) {
	// Discard any leftover cards.
	for _, hand := range g.Hands {
		g.Discard = g.Discard.PushMany(hand)
	}
	// Create new hands.
	players := len(g.Players)
	g.Hands = make([]card.Deck, players)
	per := cards.Len() / players
	for p := range g.Players {
		g.Hands[p], cards = cards.PopN(per)
		g.Hands[p] = g.Hands[p].Sort()
	}
	g.StartHand()
}

func (g *Game) StartHand() {
	g.Actions = make([]Actioner, len(g.Players))
}

func (g *Game) HasChosenAction(player int) bool {
	return g.Actions[player] != nil
}

func (g *Game) CheckHandComplete() {
	for p := range g.Players {
		if g.Actions[p] == nil || !g.Actions[p].IsComplete() {
			return
		}
	}
	for p := range g.Players {
		if pre, ok := g.Actions[p].(PreActionExecuteHandler); ok {
			pre.HandlePreActionExecute(p, g)
		}
	}
	for p := range g.Players {
		g.Actions[p].Execute(p, g)
	}
	for p := range g.Players {
		if post, ok := g.Actions[p].(PostActionExecuteHandler); ok {
			post.HandlePostActionExecute(p, g)
		}
	}
	if len(g.Hands[0]) == 1 {
		g.StartRound(g.Round + 1)
	} else {
		if g.Round%2 == 1 {
			g.Log.Add(log.NewPublicMessage("Passing hands clockwise"))
			last := len(g.Hands) - 1
			newHands := []card.Deck{g.Hands[last]}
			newHands = append(newHands, g.Hands[:last]...)
			g.Hands = newHands
		} else {
			g.Log.Add(log.NewPublicMessage("Passing hands anticlockwise"))
			newHands := append([]card.Deck{}, g.Hands[1:]...)
			newHands = append(newHands, g.Hands[0])
			g.Hands = newHands
		}
		g.StartHand()
	}
}

func (g *Game) PlayerList() []string {
	return nil
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	whose := []string{}
	for pNum, p := range g.Players {
		if g.Actions[pNum] == nil || !g.Actions[pNum].IsComplete() {
			whose = append(whose, p)
		}
	}
	return whose
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for pNum, p := range g.Players {
		if player == p {
			return pNum, true
		}
	}
	return 0, false
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func (g *Game) NumFromPlayer(player, n int) int {
	newP := (player + n) % len(g.Players)
	if newP < 0 {
		newP += len(g.Players)
	}
	return newP
}

func (g *Game) IsNeighbour(player, target int) bool {
	return g.PlayerLeft(player) == target || g.PlayerRight(player) == target
}

func (g *Game) PlayerLeft(player int) int {
	return g.NumFromPlayer(player, DirLeft)
}

func (g *Game) PlayerRight(player int) int {
	return g.NumFromPlayer(player, DirRight)
}
