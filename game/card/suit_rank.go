package card

type SuitRankCard struct {
	Suit    int
	Rank    int
	AceHigh bool
}

// Returns an ordered value corresponding to the value of the rank.  Takes into
// account if Ace is high card.
func (c SuitRankCard) RankValue() int {
	if c.Rank == STANDARD_52_RANK_ACE && c.AceHigh {
		return STANDARD_52_RANK_KING + 1
	}
	return c.Rank
}

// Sort by suit first, then card
func (c SuitRankCard) Compare(otherC Comparer) (int, bool) {
	otherCStandard, ok := otherC.(SuitRankCard)
	if !ok {
		// Different types
		return 0, false
	}
	if c.Suit < otherCStandard.Suit {
		return -1, true
	} else if c.Suit > otherCStandard.Suit {
		return 1, true
	} else if c.RankValue() < otherCStandard.RankValue() {
		return -1, true
	} else if c.RankValue() > otherCStandard.RankValue() {
		return 1, true
	}
	return 0, true
}
