package blackjack

import "github.com/Miniand/brdg.me/game/card"

type Hand struct {
	Cards       card.Deck
	Bet         int
	Insured     bool
	DoubledDown bool
	Split       bool
	Surrendered bool
}

func (h Hand) Blackjack() bool {
	return h.Cards.Len() == 2 && h.Total()[0] == 21
}

func (h Hand) Bust() bool {
	totals := h.Total()
	return totals[len(totals)-1] > 21
}

func (h Hand) Total() []int {
	totals := []int{0}
	for _, cRaw := range h.Cards {
		c := cRaw.(card.SuitRankCard)
		if c.Rank == card.STANDARD_52_RANK_ACE {
			l := len(totals)
			lastVal := totals[l-1]
			totals[l-1] += 11
			totals = append(totals, lastVal+1)
		} else {
			val := c.Rank
			if c.Rank >= card.STANDARD_52_RANK_10 {
				val = 10
			}
			for i := range totals {
				totals[i] += val
			}
		}
	}
	return totals
}
