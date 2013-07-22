package card

type SuitValueCard struct {
	Suit  int
	Value int
}

// Sort by suit first, then card
func (c SuitValueCard) Compare(otherC Comparer) (int, bool) {
	otherCStandard, ok := otherC.(SuitValueCard)
	if !ok {
		// Different types
		return 0, false
	}
	if c.Suit < otherCStandard.Suit {
		return -1, true
	} else if c.Suit > otherCStandard.Suit {
		return 1, true
	} else if c.Value < otherCStandard.Value {
		return -1, true
	} else if c.Value > otherCStandard.Value {
		return 1, true
	}
	return 0, true
}
