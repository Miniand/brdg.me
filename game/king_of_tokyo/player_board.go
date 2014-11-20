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

func (b *PlayerBoard) Things() []interface{} {
	things := []interface{}{}
	for _, c := range b.Cards {
		things = append(things, c)
	}
	for _, t := range b.Tokens {
		things = append(things, t)
	}
	return AllThings(things)
}

func AllThings(things []interface{}) []interface{} {
	allThings := append([]interface{}{}, things...)
	for _, t := range things {
		if withThings, ok := t.(HasThings); ok {
			allThings = append(allThings, AllThings(withThings.Things())...)
		}
	}
	return allThings
}
