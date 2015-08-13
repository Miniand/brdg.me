package splendor

import "github.com/Miniand/brdg.me/game/cost"

type PlayerBoard struct {
	Cards   []Card
	Reserve []Card
	Nobles  []Noble
	Tokens  cost.Cost
}

func NewPlayerBoard() PlayerBoard {
	return PlayerBoard{
		Cards:   []Card{},
		Reserve: []Card{},
		Nobles:  []Noble{},
		Tokens:  cost.Cost{},
	}
}

func (pb PlayerBoard) Bonuses() cost.Cost {
	bonuses := cost.Cost{}
	for _, c := range pb.Cards {
		bonuses[c.Resource]++
	}
	return bonuses
}

func (pb PlayerBoard) BuyingPower() cost.Cost {
	return pb.Bonuses().Add(pb.Tokens)
}

func (pb PlayerBoard) CanAfford(cost cost.Cost) bool {
	return CanAfford(pb.BuyingPower(), cost)
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
