package poker

import (
	"github.com/beefsack/brdg.me/game/card"
	"testing"
)

func buildHandByRanks(ranks []int) card.Deck {
	d := card.Deck{}
	for _, r := range ranks {
		d = d.Push(card.SuitRankCard{
			Suit:    card.STANDARD_52_SUIT_CLUBS,
			Rank:    r,
			AceHigh: true,
		})
	}
	return d
}

func TestCardsBySuit(t *testing.T) {
	hand := card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_ACE,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_8,
		},
	}
	cardsBySuit := CardsBySuit(hand)
	if len(cardsBySuit[card.STANDARD_52_SUIT_DIAMONDS]) != 3 {
		t.Fatal("Expected diamonds to be 3")
	}
	rank := cardsBySuit[card.STANDARD_52_SUIT_DIAMONDS][0].(card.SuitRankCard).Rank
	if rank != card.STANDARD_52_RANK_ACE {
		t.Fatal("Expected first diamond to be ace, got", rank)
	}
}

func TestIsStraight(t *testing.T) {
	hand := buildHandByRanks([]int{2, 6, 3, 8, 6})
	ok, _, _ := IsStraight(hand)
	if ok {
		t.Fatal("Detected as straight but isn't")
	}
	hand = buildHandByRanks([]int{2, 6, 3, 4, 5})
	ok, highCard, _ := IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 6 {
		t.Fatal("Expected high card of 6, got", highCard)
	}
	hand = buildHandByRanks([]int{2, 6, 3, 4, 5, 4})
	ok, highCard, _ = IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 6 {
		t.Fatal("Expected high card of 6, got", highCard)
	}
	// Ace as low card
	hand = buildHandByRanks([]int{2, 1, 3, 5, 4})
	ok, highCard, _ = IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 5 {
		t.Fatal("Expected high card of 5, got", highCard)
	}
	// Ace as high card
	hand = buildHandByRanks([]int{11, 10, 13, 12, 1})
	ok, highCard, _ = IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 14 {
		t.Fatal("Expected high card of 14, got", highCard)
	}
}

func TestStraightFlush(t *testing.T) {

}
