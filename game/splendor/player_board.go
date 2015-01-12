package splendor

type PlayerBoard struct {
	Cards   []Card
	Reserve []Card
	Nobles  []Noble
	Tokens  Amount
}

func NewPlayerBoard() PlayerBoard {
	return PlayerBoard{
		Cards:   []Card{},
		Reserve: []Card{},
		Nobles:  []Noble{},
		Tokens:  Amount{},
	}
}

func (pb PlayerBoard) BuyingPower() Amount {
	power := Amount{}
	for _, r := range Resources {
		power[r] = pb.Tokens[r]
	}
	for _, c := range pb.Cards {
		power[c.Resource]++
	}
	return power
}

func (pb PlayerBoard) CanAfford(amount Amount) bool {
	buyingPower := pb.BuyingPower()
	short := 0
	for r, n := range amount {
		if buyingPower[r] < n {
			short += n - buyingPower[r]
		}
	}
	return short <= buyingPower[Gold]
}
