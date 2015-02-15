package seven_wonders

type CardPlayFinalCard struct {
	Card
}

func (c CardPlayFinalCard) PlayFinalCard() bool {
	return true
}

func (c CardPlayFinalCard) SuppString() string {
	return "Can play the final card of each age"
}
