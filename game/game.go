package game

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/acquire"
	"github.com/Miniand/brdg.me/game/battleship"
	"github.com/Miniand/brdg.me/game/category_5"
	"github.com/Miniand/brdg.me/game/farkle"
	"github.com/Miniand/brdg.me/game/for_sale"
	"github.com/Miniand/brdg.me/game/liars_dice"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/game/lost_cities"
	"github.com/Miniand/brdg.me/game/modern_art"
	"github.com/Miniand/brdg.me/game/no_thanks"
	"github.com/Miniand/brdg.me/game/roll_through_the_ages"
	"github.com/Miniand/brdg.me/game/splendor"
	"github.com/Miniand/brdg.me/game/starship_catan"
	"github.com/Miniand/brdg.me/game/sushizock"
	"github.com/Miniand/brdg.me/game/texas_holdem"
	"github.com/Miniand/brdg.me/game/tic_tac_toe"
)

type Playable interface {
	Commands() []command.Command
	Name() string
	Identifier() string
	Encode() ([]byte, error)
	Decode([]byte) error
	RenderForPlayer(string) (string, error)
	Start([]string) error
	PlayerList() []string
	IsFinished() bool
	Winners() []string
	WhoseTurn() []string
	GameLog() *log.Log
}

type Eliminator interface {
	EliminatedPlayerList() []string
}

type Botter interface {
	BotPlay(player string) error
}

// The actual list of games, for a game to be active in the app it needs to be
// in here
func gameList() []Playable {
	return []Playable{
		&acquire.Game{},
		&battleship.Game{},
		&category_5.Game{},
		&farkle.Game{},
		&for_sale.Game{},
		&liars_dice.Game{},
		&lost_cities.Game{},
		&modern_art.Game{},
		&no_thanks.Game{},
		&roll_through_the_ages.Game{},
		&splendor.Game{},
		&starship_catan.Game{},
		&sushizock.Game{},
		&texas_holdem.Game{},
		&tic_tac_toe.Game{},
	}
}

// Return constructors for each game type
func Collection() map[string]func([]string) (Playable, error) {
	collection := map[string]func([]string) (Playable, error){}
	for name, raw := range RawCollection() {
		// Wrap in a function call to preserve raw for this iteration
		func(g Playable) {
			collection[name] = func(players []string) (Playable, error) {
				err := g.Start(players)
				return g, err
			}
		}(raw)
	}
	return collection
}

// Returns a collection of the raw games used for loading
func RawCollection() map[string]Playable {
	collection := map[string]Playable{}
	for _, g := range gameList() {
		collection[g.Identifier()] = g
	}
	return collection
}
