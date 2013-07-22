package card

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

type Deck []Card

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

func (d Deck) Sort() Deck {
	sort.Sort(&d)
	return d
}

func (d Deck) Len() int {
	return len(d)
}

func (d Deck) Less(i, j int) bool {
	result, comparable := d[i].Compare(d[j])
	return comparable && result < 0
}

func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Deck) Push(card Card) Deck {
	return d.PushMany([]Card{card})
}

func (d Deck) PushMany(cards []Card) Deck {
	newDeck := make(Deck, d.Len()+len(cards))
	copy(newDeck, d)
	copy(newDeck[d.Len():], cards)
	return newDeck
}

func (d Deck) Pop() (Card, Deck, error) {
	cards, newDeck, err := d.PopN(1)
	if err != nil {
		return nil, nil, err
	}
	card := cards[0]
	return card, newDeck, nil
}

func (d Deck) PopN(n int) (Deck, Deck, error) {
	if d.Len() < n {
		return nil, nil, errors.New("Not enough cards to pop")
	}
	return d[d.Len()-n:], d[:d.Len()-n], nil
}

func (d Deck) Unshift(card Card) Deck {
	return append([]Card{card}, d...)
}

func (d Deck) UnshiftMany(cards []Card) Deck {
	return append(cards, d...)
}

func (d Deck) Shift() (Card, Deck, error) {
	cards, newDeck, err := d.ShiftN(1)
	if err != nil {
		return nil, nil, err
	}
	card := cards[0]
	return card, newDeck, nil
}

func (d Deck) ShiftN(n int) (Deck, Deck, error) {
	if d.Len() < n {
		return nil, nil, errors.New("Not enough cards to shift")
	}
	return d[:n], d[n:], nil
}
