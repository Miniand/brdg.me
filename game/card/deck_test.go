package card

import (
	"fmt"
	"reflect"
	"testing"
)

// Example of manually building a deck for our own game type.
func ExampleDeckBuild() {
	deck := Deck{}
	for i := 1; i <= 3; i++ {
		deck = deck.Push(SuitRankCard{
			Suit: STANDARD_52_SUIT_HEARTS,
			Rank: i,
		})
	}
	fmt.Printf("The first card rank is %d\n", deck[0].(SuitRankCard).Rank)
	fmt.Printf("The second card rank is %d\n", deck[1].(SuitRankCard).Rank)
	fmt.Printf("The third card rank is %d\n", deck[2].(SuitRankCard).Rank)
	// Output:
	// The first card rank is 1
	// The second card rank is 2
	// The third card rank is 3
}

// Example of shuffling a normal deck and dealing 5 cards each to three players.
// We make sure the player hands remain sorted.
func ExampleDeal() {
	deck := Standard52Deck()
	deck = deck.Shuffle()
	player1Hand, deck := deck.PopN(5)
	player1Hand = player1Hand.Sort()
	player2Hand, deck := deck.PopN(5)
	player2Hand = player2Hand.Sort()
	player3Hand, deck := deck.PopN(5)
	player3Hand = player3Hand.Sort()
	fmt.Printf("Player 1 hand size is %d\n", player1Hand.Len())
	fmt.Printf("Player 2 hand size is %d\n", player2Hand.Len())
	fmt.Printf("Player 3 hand size is %d\n", player3Hand.Len())
	fmt.Printf("There are %d cards remaining in the deck", deck.Len())
	// Output:
	// Player 1 hand size is 5
	// Player 2 hand size is 5
	// Player 3 hand size is 5
	// There are 37 cards remaining in the deck
}

// Example of examining ranks in the card struct, we need to type assert it
// first.  Shift takes a card from the front instead of the back.  We will also
// peek the next card on the deck without actually taking it.
func ExampleExamineCard() {
	deck := Standard52Deck() // We won't shuffle it so we know the order
	player1Card, deck := deck.Shift()
	player2Card, deck := deck.Shift()
	nextCard, _ := deck.Shift()
	fmt.Printf("Player 1's card rank is %d\n",
		player1Card.(SuitRankCard).Rank)
	fmt.Printf("Player 2's card rank is %d\n",
		player2Card.(SuitRankCard).Rank)
	fmt.Printf("The next card rank is %d\n", nextCard.(SuitRankCard).Rank)
	fmt.Printf("The remaining card count is still %d\n", deck.Len())
	// Output:
	// Player 1's card rank is 1
	// Player 2's card rank is 2
	// The next card rank is 3
	// The remaining card count is still 50
}

func testIsStandardDeck(d Deck, t *testing.T) {
	if d.Len() != 52 {
		t.Fatal("Deck is not 52 cards")
	}
	sortedD := d.Sort()
	i := 0
	for suit := STANDARD_52_SUIT_SPADES; suit <= STANDARD_52_SUIT_CLUBS; suit++ {
		for rank := STANDARD_52_RANK_ACE; rank <= STANDARD_52_RANK_KING; rank++ {
			result, comparable := sortedD[i].Compare(SuitRankCard{
				Suit: suit,
				Rank: rank,
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
	if reflect.DeepEqual(d, d.Shuffle()) {
		t.Fatal("was the same after shuffling")
	}
}

func TestSort(t *testing.T) {
	d := Standard52Deck().Shuffle()
	dClone := append(Deck{}, d...)
	newD := d.Sort()
	if !reflect.DeepEqual(newD, Standard52Deck()) {
		t.Fatal("newD is not sorted")
	}
	if !reflect.DeepEqual(d, dClone) {
		t.Fatal("d was not preserved")
	}
}

func TestPush(t *testing.T) {
	d := Standard52Deck()
	c := SuitRankCard{
		Suit: 50,
		Rank: 50,
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
		SuitRankCard{
			Suit: 50,
			Rank: 50,
		},
		SuitRankCard{
			Suit: 51,
			Rank: 51,
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
	shouldBeCard := SuitRankCard{
		Suit: STANDARD_52_SUIT_SPADES,
		Rank: STANDARD_52_RANK_KING,
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
	result, _ := cards[0].Compare(SuitRankCard{
		Suit: STANDARD_52_SUIT_SPADES,
		Rank: STANDARD_52_RANK_QUEEN,
	})
	if result != 0 {
		t.Fatal("First card popped wasn't Queen of Spades")
	}
	result, _ = cards[1].Compare(SuitRankCard{
		Suit: STANDARD_52_SUIT_SPADES,
		Rank: STANDARD_52_RANK_KING,
	})
	if result != 0 {
		t.Fatal("Second card popped wasn't King of Spades")
	}
}

func TestUnshift(t *testing.T) {
	d := Standard52Deck()
	c := SuitRankCard{
		Suit: 50,
		Rank: 50,
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
		SuitRankCard{
			Suit: 50,
			Rank: 50,
		},
		SuitRankCard{
			Suit: 51,
			Rank: 51,
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
	shouldBeCard := SuitRankCard{
		Suit: STANDARD_52_SUIT_CLUBS,
		Rank: STANDARD_52_RANK_ACE,
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
	result, _ := cards[0].Compare(SuitRankCard{
		Suit: STANDARD_52_SUIT_CLUBS,
		Rank: STANDARD_52_RANK_ACE,
	})
	if result != 0 {
		t.Fatal("First card shifted wasn't Ace of Clubs")
	}
	result, _ = cards[1].Compare(SuitRankCard{
		Suit: STANDARD_52_SUIT_CLUBS,
		Rank: STANDARD_52_RANK_2,
	})
	if result != 0 {
		t.Fatal("Second card shifted wasn't Two of Clubs")
	}
}

func TestToSuitRankCards(t *testing.T) {
	d := Standard52Deck()
	cards := d.ToSuitRankCards()
	if d.Len() != len(cards) {
		t.Fatal("Length of deck doesn't match length of cards array")
	}
}
