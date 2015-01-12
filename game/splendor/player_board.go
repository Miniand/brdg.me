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

func (pb PlayerBoard) Bonuses() Amount {
	bonuses := Amount{}
	for _, c := range pb.Cards {
		bonuses[c.Resource]++
	}
	return bonuses
}

func (pb PlayerBoard) BuyingPower() Amount {
	return pb.Bonuses().Add(pb.Tokens)
}

func (pb PlayerBoard) CanAfford(cost Amount) bool {
	return pb.BuyingPower().CanAfford(cost)
}
