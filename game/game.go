package game

import (
	"github.com/beefsack/boredga.me/game/tic_tac_toe"
)

type Playable interface {
	PlayerAction(string, string, []string) error
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
}

// The actual list of games, for a game to be active in the app it needs to be
// in here
func gameList() []Playable {
	return []Playable{
		&tic_tac_toe.Game{},
	}
}

// Return constructors for each game type
func Collection() map[string]func([]string) (Playable, error) {
	collection := map[string]func([]string) (Playable, error){}
	for name, raw := range RawCollection() {
		collection[name] = func(players []string) (Playable, error) {
			err := raw.Start(players)
			return raw, err
		}
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
