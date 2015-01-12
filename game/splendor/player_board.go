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

func (pb PlayerBoard) Prestige() int {
	prestige := 0
	for _, c := range pb.Cards {
		prestige += c.Prestige
	}
	for _, n := range pb.Nobles {
		prestige += n.Prestige
	}
	return prestige
}

func (pb *PlayerBoard) Pay(amount Amount) {
	offset := pb.Bonuses().Subtract(amount)
	for _, g := range Gems {
		if offset[g] < 0 {
			// We didn't have enough just with bonuses
			pb.Tokens[g] += offset[g]
			if pb.Tokens[g] < 0 {
				// We didn't have enough normal tokens either, use gold
				pb.Tokens[Gold] += pb.Tokens[g]
				pb.Tokens[g] = 0
			}
		}
	}
}
