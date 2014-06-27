package starship_catan

import "errors"

type Game struct {
	Players      []string
	PlayerBoards [2]*PlayerBoard
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("this game requires two players")
	}
	g.Players = players
	return nil
}
