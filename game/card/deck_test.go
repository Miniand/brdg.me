package card

import (
	"testing"
)

func testIsStandardDeck(d Deck, t *testing.T) {
	if d.Len() != 52 {
		t.Fatal("Deck is not 52 cards")
	}
	sortedD := d.Sort()
	i := 0
	for suit := STANDARD_52_SUIT_SPADES; suit <= STANDARD_52_SUIT_CLUBS; suit++ {
		for value := STANDARD_52_VALUE_ACE; value <= STANDARD_52_VALUE_KING; value++ {
			result, comparable := sortedD[i].Compare(SuitValueCard{
				Suit:  suit,
				Value: value,
			})
			if !comparable || result != 0 {
				t.Fatal("Deck is not standard, card", sortedD[i])
			}
			i++
		}
	}
}

func TestContains(t *testing.T) {
	d := Standard52Deck()
	c := d[3]
	d = d.Push(c)
	if d.Contains(c) != 2 {
		t.Fatal("Did not count 2 cards")
	}
}

func TestRemove(t *testing.T) {
	d := Standard52Deck()
	c := d[3]
	newD, removed := d.Remove(c, -1)
	testIsStandardDeck(d, t)
	if removed != 1 {
		t.Fatal("Didn't report 1 card was removed")
	}
	if newD.Len() != 51 {
		t.Fatal("Length of deck is not 51")
	}
	// Try again, 0 should be removed
	newD, removed = newD.Remove(c, -1)
	if removed != 0 {
		t.Fatal("Didn't report 0 cards were removed")
	}
	if newD.Len() != 51 {
		t.Fatal("Length of deck is not 51")
	}
}

func TestShuffle(t *testing.T) {
	d := Standard52Deck()
	newD := d.Shuffle()
	result, _ := newD[0].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_ACE,
	})
	if result == 0 {
		t.Fatal("Ace of Spades didn't move from the first spot, though there's a chance it was still shuffled")
	}
}

func TestSort(t *testing.T) {
	d := Standard52Deck().Shuffle()
	newD := d.Sort()
	result, _ := d[0].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_ACE,
	})
	if result == 0 {
		t.Fatal("Original neck is no longer shuffled")
	}
	result, _ = newD[0].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_ACE,
	})
	if result != 0 {
		t.Fatal("New deck did not sort to put the ace of clubs first")
	}
}

func TestPush(t *testing.T) {
	d := Standard52Deck()
	c := SuitValueCard{
		Suit:  50,
		Value: 50,
	}
	newD := d.Push(c)
	if len(d) != 52 {
		t.Fatal("Push modified original deck")
	}
	if len(newD) != 53 {
		t.Fatal("Deck is not 53 cards")
	}
	result, _ := newD[52].Compare(c)
	if result != 0 {
		t.Fatal("Pushed card is not last card")
	}
}

func TestPushMany(t *testing.T) {
	d := Standard52Deck()
	cards := []Card{
		SuitValueCard{
			Suit:  50,
			Value: 50,
		},
		SuitValueCard{
			Suit:  51,
			Value: 51,
		},
	}
	newD := d.PushMany(cards)
	if len(d) != 52 {
		t.Fatal("Push modified original deck")
	}
	if len(newD) != 54 {
		t.Fatal("Deck is not 54 cards")
	}
	result, _ := newD[52].Compare(cards[0])
	if result != 0 {
		t.Fatal("First pushed card is not second last card")
	}
	result, _ = newD[53].Compare(cards[1])
	if result != 0 {
		t.Fatal("Second pushed card is not last card")
	}
}

func TestPop(t *testing.T) {
	d := Standard52Deck()
	c, newD := d.Pop()
	if len(d) != 52 {
		t.Fatal("Pop modified original deck")
	}
	if len(newD) != 51 {
		t.Fatal("Deck is not 51 cards")
	}
	shouldBeCard := SuitValueCard{
		Suit:  STANDARD_52_SUIT_SPADES,
		Value: STANDARD_52_VALUE_KING,
	}
	result, _ := c.Compare(shouldBeCard)
	if result != 0 {
		t.Fatal("Card popped was not", shouldBeCard, ", got:", c)
	}
}

func TestPopN(t *testing.T) {
	d := Standard52Deck()
	cards, newD := d.PopN(2)
	if len(d) != 52 {
		t.Fatal("PopN modified original deck")
	}
	if len(newD) != 50 {
		t.Fatal("Deck is not 50 cards")
	}
	if len(cards) != 2 {
		t.Fatal("Taken cards isn't length 2")
	}
	result, _ := cards[0].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_SPADES,
		Value: STANDARD_52_VALUE_QUEEN,
	})
	if result != 0 {
		t.Fatal("First card popped wasn't Queen of Spades")
	}
	result, _ = cards[1].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_SPADES,
		Value: STANDARD_52_VALUE_KING,
	})
	if result != 0 {
		t.Fatal("Second card popped wasn't King of Spades")
	}
}

func TestUnshift(t *testing.T) {
	d := Standard52Deck()
	c := SuitValueCard{
		Suit:  50,
		Value: 50,
	}
	newD := d.Unshift(c)
	if len(d) != 52 {
		t.Fatal("Unshift modified original deck")
	}
	if len(newD) != 53 {
		t.Fatal("Deck is not 53 cards")
	}
	result, _ := newD[0].Compare(c)
	if result != 0 {
		t.Fatal("Unshifted card is not first card")
	}
}

func TestUnshiftMany(t *testing.T) {
	d := Standard52Deck()
	cards := []Card{
		SuitValueCard{
			Suit:  50,
			Value: 50,
		},
		SuitValueCard{
			Suit:  51,
			Value: 51,
		},
	}
	newD := d.UnshiftMany(cards)
	if len(d) != 52 {
		t.Fatal("Unshift modified original deck")
	}
	if len(newD) != 54 {
		t.Fatal("Deck is not 54 cards")
	}
	result, _ := newD[0].Compare(cards[0])
	if result != 0 {
		t.Fatal("First unshifted card is not first card")
	}
	result, _ = newD[1].Compare(cards[1])
	if result != 0 {
		t.Fatal("Second unshifted card is not second card")
	}
}

func TestShift(t *testing.T) {
	d := Standard52Deck()
	c, newD := d.Shift()
	if len(d) != 52 {
		t.Fatal("Shift modified original deck")
	}
	if len(newD) != 51 {
		t.Fatal("Deck is not 51 cards")
	}
	shouldBeCard := SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_ACE,
	}
	result, _ := c.Compare(shouldBeCard)
	if result != 0 {
		t.Fatal("Card shifted was not", shouldBeCard, ", got:", c)
	}
}

func TestShiftN(t *testing.T) {
	d := Standard52Deck()
	cards, newD := d.ShiftN(2)
	if len(d) != 52 {
		t.Fatal("ShiftN modified original deck")
	}
	if len(newD) != 50 {
		t.Fatal("Deck is not 50 cards")
	}
	if len(cards) != 2 {
		t.Fatal("Taken cards isn't length 2")
	}
	result, _ := cards[0].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_ACE,
	})
	if result != 0 {
		t.Fatal("First card shifted wasn't Ace of Clubs")
	}
	result, _ = cards[1].Compare(SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_2,
	})
	if result != 0 {
		t.Fatal("Second card shifted wasn't Two of Clubs")
	}
}

func TestToSuitValueCards(t *testing.T) {
	d := Standard52Deck()
	cards := d.ToSuitValueCards()
	if d.Len() != len(cards) {
		t.Fatal("Length of deck doesn't match length of cards array")
	}
}
