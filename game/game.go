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
	RenderForPlayer(string) (error, string)
	Start([]string) error
}

// The actual list of games, for a game to be active in the app it needs to be
// in here
func gameList() []Playable {
	return []Playable{
		&tic_tac_toe.Game{},
	}
}

// Return constructors for each game type
func Collection() map[string]func([]string) (error, Playable) {
	collection := map[string]func([]string) (error, Playable){}
	for name, raw := range RawCollection() {
		collection[name] = func(players []string) (error, Playable) {
			err := raw.Init(players)
			return err, raw
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
