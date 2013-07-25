package card

type SuitRankCard struct {
	Suit int
	Rank int
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
	} else if c.Rank < otherCStandard.Rank {
		return -1, true
	} else if c.Rank > otherCStandard.Rank {
		return 1, true
	}
	return 0, true
}
