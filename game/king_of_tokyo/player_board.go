package king_of_tokyo

type PlayerBoard struct {
	Health int
	VP     int
	Energy int
	Cards  []CardBase
	Tokens []interface{}
}

func NewPlayerBoard() *PlayerBoard {
	return &PlayerBoard{
		Health: 10,
		Cards:  []CardBase{},
		Tokens: []interface{}{},
	}
}
