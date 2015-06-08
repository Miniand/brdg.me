package seven_wonders

import "encoding/gob"

func init() {
	gob.Register(CardPlayFinalCard{})
}

type CardPlayFinalCard struct {
	Card
}

func (c CardPlayFinalCard) PlayFinalCard() bool {
	return true
}

func (c CardPlayFinalCard) SuppString() string {
	return "Can play the final card of each age"
}
