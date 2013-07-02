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
}

// Return constructors for each game type
func Collection() map[string]func([]string) (error, Playable) {
	return map[string]func([]string) (error, Playable){
		tic_tac_toe.RawGame().Identifier(): func(players []string) (error, Playable) {
			return tic_tac_toe.NewGame(players)
		},
	}
}

// Returns a collection of the raw games used for loading
func RawCollection() map[string]Playable {
	return map[string]Playable{
		tic_tac_toe.RawGame().Identifier(): tic_tac_toe.RawGame(),
	}
}
