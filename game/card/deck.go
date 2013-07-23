

package card

import (
	"math/rand"
	"sort"
	"time"
)

type Deck []Card

// Counts how many instances of a card appears in a deck
func (d Deck) Contains(card Card) int {
	_, count := d.Remove(card, -1)
	return count
}

// Removes cards from the deck.  Removes up to n cards, or removes all instances
// if n = -1.  Returns deck with cards removed and the cound of cards removed.
func (d Deck) Remove(removeCard Card, n int) (Deck, int) {
	removed := 0
	newD := d[:0]
	for _, c := range d {
		result, comparable := c.Compare(removeCard)
		if comparable && result == 0 && (n < 0 || removed < n) {
			removed++
		} else {
			newD = newD.Push(c)
		}
	}
	return newD, removed
}

// Returns a shuffled version of the deck
func (d Deck) Shuffle() Deck {
	l := d.Len()
	if l <= 1 {
		return d
	}
	newD := Deck{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(l)
	for i := 0; i < l; i++ {
		newD = newD.Push(d[perm[i]])
	}
	return newD
}

// Returns a sorted version of the deck
func (d Deck) Sort() Deck {
	newD := make(Deck, d.Len())
	copy(newD, d)
	sort.Sort(newD)
	return newD
}

// Returns the length of the deck
func (d Deck) Len() int {
	return len(d)
}

// Returns whether the item at offset i is less than the item at offset j, used
// for sorting
func (d Deck) Less(i, j int) bool {
	result, comparable := d[i].Compare(d[j])
	return comparable && result < 0
}

// Swaps the cards at indexes i and j
func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Returns a new deck with the card appended to the end
func (d Deck) Push(card Card) Deck {
	return d.PushMany([]Card{card})
}

// Returns a new deck with the cards appended to the end
func (d Deck) PushMany(cards []Card) Deck {
	newDeck := make(Deck, d.Len()+len(cards))
	copy(newDeck, d)
	copy(newDeck[d.Len():], cards)
	return newDeck
}

// Returns the last card in the deck, and the remaining cards
func (d Deck) Pop() (Card, Deck) {
	cards, newDeck := d.PopN(1)
	card := cards[0]
	return card, newDeck
}

// Returns the last n cards in the deck, and the remaining cards
func (d Deck) PopN(n int) (Deck, Deck) {
	if d.Len() < n {
		panic("Not enough cards to pop")
	}
	return d[d.Len()-n:], d[:d.Len()-n]
}

// Returns a new deck with the card prepended to the front
func (d Deck) Unshift(card Card) Deck {
	return append([]Card{card}, d...)
}

// Returns a new deck with the cards prepended to the front
func (d Deck) UnshiftMany(cards []Card) Deck {
	return append(cards, d...)
}

// Returns the first card in the deck, and the remaining cards
func (d Deck) Shift() (Card, Deck) {
	cards, newDeck := d.ShiftN(1)
	card := cards[0]
	return card, newDeck
}

// Returns the first n cards in the deck, and the remaining cards
func (d Deck) ShiftN(n int) (Deck, Deck) {
	if d.Len() < n {
		panic("Not enough cards to shift")
	}
	return d[:n], d[n:]
}

// Converts deck to a flat array of SuitValueCard
func (d Deck) ToSuitValueCards() []SuitValueCard {
	cards := make([]SuitValueCard, d.Len())
	for i, c := range d {
		cards[i] = c.(SuitValueCard)
	}
	return cards
}
