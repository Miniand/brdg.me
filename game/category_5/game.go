package category_5

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players      []string
	Deck         []Card
	Discard      []Card
	Points       map[int]int
	Hands        map[int][]Card
	PlayerCards  map[int][]Card
	Plays        map[int]Card
	Board        [4][]Card
	Resolving    bool
	ChoosePlayer int
	Log          *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
		ChooseCommand{},
	}
}

func (g *Game) Name() string {
	return "Category 5"
}

func (g *Game) Identifier() string {
	return "category_5"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 || len(players) > 10 {
		return errors.New("this game is for 2-10 players")
	}
	g.Log = log.New()
	g.Players = players
	g.Deck = Shuffle(Deck())
	g.Discard = []Card{}
	g.Points = map[int]int{}
	g.Hands = map[int][]Card{}
	g.PlayerCards = map[int][]Card{}
	g.Plays = map[int]Card{}
	g.Board = [4][]Card{{}, {}, {}, {}}
	for p, _ := range g.Players {
		g.Hands[p] = []Card{}
		g.PlayerCards[p] = []Card{}
	}
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	// Discard cards on the board
	for i, b := range g.Board {
		g.DiscardCards(b)
		g.Board[i] = []Card{}
	}
	// Discard the player cards
	for p, _ := range g.Players {
		g.DiscardCards(g.PlayerCards[p])
		g.PlayerCards[p] = []Card{}
	}
	// Start each row with a card
	for i, _ := range g.Board {
		g.Board[i] = append(g.Board[i], g.DrawCards(1)...)
	}
	// Each player gets 10 cards
	for p, _ := range g.Players {
		g.Hands[p] = SortCards(g.DrawCards(10))
	}
	g.Log.Add(log.NewPublicMessage(
		"Starting a new round, dealing 10 cards to each player"))
}

func (g *Game) ResolvePlays() {
	g.Resolving = true
	for {
		// Find who has the next lowest card
		lowestCard := Card(0)
		lowestPlayer := 0
		for p, _ := range g.Players {
			if g.Plays[p] == 0 {
				continue
			}
			if lowestCard == 0 || g.Plays[p] < lowestCard {
				lowestCard = g.Plays[p]
				lowestPlayer = p
			}
		}
		if lowestCard == 0 {
			// None left, we've resolved all
			break
		}
		// Find which row it goes in
		closestCard := Card(0)
		closestRow := 0
		for i, row := range g.Board {
			lastCard := row[len(row)-1]
			if lastCard < lowestCard && (closestCard == 0 || lastCard > closestCard) {
				closestCard = lastCard
				closestRow = i
			}
		}
		if closestCard == 0 {
			// The card is lower than all rows, player gets to choose row
			g.ChoosePlayer = lowestPlayer
			return
		} else if len(g.Board[closestRow]) == 5 {
			// Row is full, gotta take it
			g.PlayerCards[lowestPlayer] = append(
				g.PlayerCards[lowestPlayer], g.Board[closestRow]...)
			g.Board[closestRow] = []Card{lowestCard}
		} else {
			// Just slot the card into the row
			g.Board[closestRow] = append(g.Board[closestRow], lowestCard)
		}
		g.Plays[lowestPlayer] = 0
	}
	g.Resolving = false
	if len(g.Hands[0]) == 0 {
		g.EndRound()
	}
}

func (g *Game) EndRound() {
	for p, _ := range g.Players {
		total := 0
		for _, c := range g.PlayerCards[p] {
			total += int(c)
		}
		g.Points[p] += total
	}
	if !g.IsFinished() {
		g.StartRound()
	}
}

func (g *Game) DrawCards(n int) []Card {
	cards := []Card{}
	if l := len(g.Deck); l >= n {
		cards, g.Deck = TakeCards(g.Deck, n)
	} else {
		cards = append(cards, g.Deck...)
		g.Deck = Shuffle(g.Discard)
		g.Discard = []Card{}
		cards = append(cards, g.DrawCards(n-l)...)
	}
	return cards
}

func (g *Game) DiscardCards(cards []Card) {
	g.Discard = append(g.Discard, cards...)
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	highestScore := 0
	for p, _ := range g.Players {
		if g.Points[p] > highestScore {
			highestScore = g.Points[p]
		}
		if len(g.Hands[p]) > 0 {
			return false
		}
	}
	return highestScore >= 66
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	lowestScore := -1
	winners := []string{}
	for p, pName := range g.Players {
		if lowestScore == -1 || g.Points[p] < lowestScore {
			lowestScore = g.Points[p]
			winners = []string{}
		}
		if g.Points[p] == lowestScore {
			winners = append(winners, pName)
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	if g.Resolving {
		return []string{g.Players[g.ChoosePlayer]}
	}
	whose := []string{}
	for p, pName := range g.Players {
		if g.Plays[p] == 0 {
			whose = append(whose, pName)
		}
	}
	return whose
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, error) {
	for pNum, p := range g.Players {
		if p == player {
			return pNum, nil
		}
	}
	return 0, errors.New("could not find player")
}
