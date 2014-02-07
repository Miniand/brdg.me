package modern_art

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	INITIAL_MONEY   = 100
	SUIT_LITE_METAL = iota
	SUIT_YOKO
	SUIT_CHRISTINE_P
	SUIT_KARL_GLITTER
	SUIT_KRYPTO
	RANK_OPEN = iota
	RANK_FIXED_PRICE
	RANK_SEALED
	RANK_DOUBLE
	RANK_ONCE_AROUND
)

var roundCards = map[int]map[int]int{
	3: map[int]int{
		0: 10,
		1: 6,
		2: 6,
		3: 0,
	},
	4: map[int]int{
		0: 9,
		1: 4,
		2: 4,
		3: 0,
	},
	5: map[int]int{
		0: 8,
		1: 3,
		2: 3,
		3: 0,
	},
}

var suits = []int{
	SUIT_LITE_METAL,
	SUIT_YOKO,
	SUIT_CHRISTINE_P,
	SUIT_KARL_GLITTER,
	SUIT_KRYPTO,
}

var suitNames = map[int]string{
	SUIT_LITE_METAL:   "Lite Metal",
	SUIT_YOKO:         "Lite Metal",
	SUIT_CHRISTINE_P:  "Lite Metal",
	SUIT_KARL_GLITTER: "Lite Metal",
	SUIT_KRYPTO:       "Lite Metal",
}

var ranks = []int{
	RANK_OPEN,
	RANK_FIXED_PRICE,
	RANK_SEALED,
	RANK_DOUBLE,
	RANK_ONCE_AROUND,
}

var rankNames = map[int]string{
	RANK_OPEN:        "Open",
	RANK_FIXED_PRICE: "Fixed Price",
	RANK_SEALED:      "Sealed",
	RANK_DOUBLE:      "Double",
	RANK_ONCE_AROUND: "Once Around",
}

var cardDistribution = map[int]map[int]int{
	SUIT_LITE_METAL: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 2,
		RANK_SEALED:      2,
		RANK_DOUBLE:      2,
		RANK_ONCE_AROUND: 3,
	},
	SUIT_YOKO: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      2,
		RANK_ONCE_AROUND: 2,
	},
	SUIT_CHRISTINE_P: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      2,
		RANK_ONCE_AROUND: 3,
	},
	SUIT_KARL_GLITTER: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      3,
		RANK_ONCE_AROUND: 3,
	},
	SUIT_KRYPTO: map[int]int{
		RANK_OPEN:        4,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      3,
		RANK_ONCE_AROUND: 3,
	},
}

type Game struct {
	Players     []string
	PlayerMoney map[int]int
	PlayerHands map[int]card.Deck
	Round       int
	Deck        card.Deck
	Log         *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Modern Art"
}

func (g *Game) Identifier() string {
	return "modern_art"
}

func RegisterGobTypes() {
	gob.Register(card.SuitRankCard{})
}

func (g *Game) Encode() ([]byte, error) {
	RegisterGobTypes()
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(g)
	return buf.Bytes(), err
}

func (g *Game) Decode(data []byte) error {
	RegisterGobTypes()
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(g)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) Start(players []string) error {
	if len(players) < 3 || len(players) > 5 {
		return errors.New("Modern Art requires between 3 and 5 players")
	}
	g.Players = players
	g.PlayerMoney = map[int]int{}
	g.PlayerHands = map[int]card.Deck{}
	for i, _ := range g.Players {
		g.PlayerMoney[i] = INITIAL_MONEY
		g.PlayerHands[i] = card.Deck{}
	}
	g.Deck = Deck().Shuffle()
	g.Log = &log.Log{}
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	numCards := roundCards[len(g.Players)][g.Round]
	if numCards <= 0 {
		return
	}
	for i, _ := range g.Players {
		cards, remaining := g.Deck.PopN(numCards)
		g.PlayerHands[i].PushMany(cards)
		g.Deck = remaining
	}
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

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func Deck() card.Deck {
	d := card.Deck{}
	for suit, suitCards := range cardDistribution {
		for rank, n := range suitCards {
			for i := 0; i < n; i++ {
				d = d.Push(card.SuitRankCard{suit, rank})
			}
		}
	}
	return d
}
