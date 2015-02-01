package alhambra

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	PhaseAction = iota
	PhasePlace
	PhaseFinalPlace
	PhaseEnd
)

var RoundScores = map[int][]int{
	TileTypePavillion: []int{1, 8, 16},
	TileTypeSeraglio:  []int{2, 9, 17},
	TileTypeArcades:   []int{3, 10, 18},
	TileTypeChambers:  []int{4, 11, 19},
	TileTypeGarden:    []int{5, 12, 20},
	TileTypeTower:     []int{6, 13, 21},
}

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
		TakeCommand{},
		PlaceCommand{},
		SwapCommand{},
		RemoveCommand{},
		DoneCommand{},
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
	gob.Register(ScoringCard{})
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

	// Inject scoring cards
	subSize := len(g.CardPile) / 5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, pos := range []int{
		subSize + r.Int()%subSize,       // In the 2nd 5th of cards
		3*subSize + r.Int()%subSize + 1, // In the 4th 5th of cards
	} {
		h2 := append(card.Deck{}, g.CardPile[pos:]...)
		g.CardPile = append(card.Deck{}, g.CardPile[:pos]...)
		g.CardPile = append(g.CardPile, ScoringCard{})
		g.CardPile = append(g.CardPile, h2...)
	}

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
	output := bytes.NewBufferString(fmt.Sprintf(
		"{{b}}It is now scoring round %d{{_b}}",
		g.Round,
	))
	for _, t := range ScoringTileTypes {
		output.WriteString(fmt.Sprintf(
			"\n{{b}}Scoring %s{{_b}}",
			RenderTileAbbr(t),
		))
		for _, rs := range g.ScoreType(t, g.Round) {
			playerStrs := []string{}
			for _, p := range rs.Players {
				g.Boards[p].Points += rs.Points
				playerStrs = append(playerStrs, g.PlayerName(p))
			}
			output.WriteString(fmt.Sprintf(
				"\n%s scored %d for having %d",
				render.CommaList(playerStrs),
				rs.Points,
				rs.TileCount,
			))
		}
	}
	output.WriteString("\n{{b}}Scoring walls{{_b}}")
	for p := range g.Players {
		wall := g.Boards[p].Grid.LongestExtWall()
		g.Boards[p].Points += wall
		output.WriteString(fmt.Sprintf(
			"\n%s scored %d for their wall",
			g.PlayerName(p),
			wall,
		))
	}
	if g.Round < 3 {
		g.Round++
	}
	g.Log.Add(log.NewPublicMessage(output.String()))
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Phase == PhaseEnd
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winners := []string{}
	score := 0
	for p, name := range g.Players {
		if g.Boards[p].Points > score {
			winners = []string{}
			score = g.Boards[p].Points
		}
		if g.Boards[p].Points == score {
			winners = append(winners, name)
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
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
	case PhaseFinalPlace:
		nextPlayer := g.CurrentPlayer
		for {
			nextPlayer = (nextPlayer + 1) % len(g.Players)
			if nextPlayer == g.CurrentPlayer {
				// Everyone has placed, final scoring
				g.Phase = PhaseEnd
				g.ScoreRound()
				break
			}
			if len(NotEmpty(g.Boards[nextPlayer].Place)) > 0 {
				g.CurrentPlayer = nextPlayer
				break
			}
		}
	}
}

func (g *Game) NextPlayer() {
	// Clean up existing turn
	reserved := []string{}
	for _, t := range g.Boards[g.CurrentPlayer].Place {
		if t.Type != TileTypeEmpty {
			reserved = append(reserved, RenderTileAbbr(t.Type))
			g.Boards[g.CurrentPlayer].Reserve = append(
				g.Boards[g.CurrentPlayer].Reserve, t)
		}
	}
	if len(reserved) > 0 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s added %s to their reserve",
			g.PlayerName(g.CurrentPlayer),
			render.CommaList(reserved),
		)))
	}
	g.Boards[g.CurrentPlayer].Place = []Tile{}
	for i, t := range g.Tiles {
		if t.Type == TileTypeEmpty {
			if len(g.TileBag) > 0 {
				g.Tiles[i] = g.TileBag[0]
				g.TileBag = g.TileBag[1:]
			} else {
				// End of the game
			}
		}
	}
	if l := g.Cards.Len(); l < 4 {
		g.Cards = g.Cards.PushMany(g.DrawCards(4 - l))
	}
	// Move to next player
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

type RoundTypeScore struct {
	Players           []int
	TileCount, Points int
}

func (g *Game) ScoreType(tileType, round int) []RoundTypeScore {
	// Group players by tile count
	byCount := map[int][]int{}
	counts := []int{}
	for p, _ := range g.Players {
		count := g.Boards[p].TileCounts()[tileType]
		if count == 0 {
			continue
		}
		if byCount[count] == nil {
			byCount[count] = []int{}
			counts = append(counts, count)
		}
		byCount[count] = append(byCount[count], p)
	}

	// Loop through counts and assign points
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	rewards := RoundScores[tileType][:round]
	scores := []RoundTypeScore{}
	for len(rewards) > 0 && len(counts) > 0 {
		rts := RoundTypeScore{
			Players: byCount[counts[0]],
		}

		n := len(rts.Players)
		if l := len(rewards); n > l {
			n = l
		}
		points := 0
		for _, r := range rewards[len(rewards)-n:] {
			points += r
		}
		rts.TileCount = counts[0]
		rts.Points = points / len(rts.Players)

		scores = append(scores, rts)
		rewards = rewards[:len(rewards)-n]
		counts = counts[1:]
	}
	return scores
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
