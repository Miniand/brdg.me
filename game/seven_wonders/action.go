package seven_wonders

import "encoding/gob"

func init() {
	gob.Register(&BuildAction{})
	gob.Register(&WonderAction{})
	gob.Register(&DiscardAction{})
}

type Actioner interface {
	IsComplete() bool
	Execute(player int, g *Game)
	Output(player int, g *Game) string
}

type DealOptioner interface {
	DealOptions(player int, g *Game) []map[int]int
	ChooseDeal(player int, g *Game, n int) error
}

type WonderAction struct {
	Card   int
	Deal   int
	Chosen bool
}

func (a WonderAction) IsComplete() bool {
	return a.Chosen
}

func (a WonderAction) DealOptions(player int, g *Game) []map[int]int {
	return nil
}
