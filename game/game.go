package game

import (
	"github.com/beefsack/boredga.me/game/tic_tac_toe"
)

type Playable interface {
	PlayerAction(string, string, []string) error
}

func Collection() map[string]func([]string) (error, Playable) {
	return map[string]func([]string) (error, Playable){
		"Tic-tac-toe": func(players []string) (error, Playable) {
			return tic_tac_toe.NewGame(players)
		},
	}
}
