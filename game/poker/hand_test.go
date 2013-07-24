package poker

import (
	"github.com/beefsack/brdg.me/game/card"
	"testing"
)

func TestSortRanks(t *testing.T) {
	ranks := []int{
		card.STANDARD_52_RANK_5,
		card.STANDARD_52_RANK_ACE,
		card.STANDARD_52_RANK_7,
		card.STANDARD_52_RANK_2,
		card.STANDARD_52_RANK_ACE,
		card.STANDARD_52_RANK_KING,
	}
	sorted := SortRanks(ranks)
	if len(sorted) != 6 {
		t.Fatal("Expected 6 cards")
	}
	if sorted[0] != card.STANDARD_52_RANK_ACE {
		t.Fatal("Expected first card to be the ace")
	}
	if sorted[1] != card.STANDARD_52_RANK_ACE {
		t.Fatal("Expected second card to be the ace")
	}
	if sorted[2] != card.STANDARD_52_RANK_KING {
		t.Fatal("Expected third card to be the king")
	}
	if sorted[3] != card.STANDARD_52_RANK_7 {
		t.Fatal("Expected fourth card to be the seven")
	}
	if sorted[4] != card.STANDARD_52_RANK_5 {
		t.Fatal("Expected fifth card to be the five")
	}
	if sorted[5] != card.STANDARD_52_RANK_2 {
		t.Fatal("Expected sixth card to be the two")
	}
}

func TestRanksBySuit(t *testing.T) {
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
	ranksBySuit := RanksBySuit(hand)
	if len(ranksBySuit[card.STANDARD_52_SUIT_DIAMONDS]) != 3 {
		t.Fatal("Expected diamonds to be 3")
	}
	if ranksBySuit[card.STANDARD_52_SUIT_DIAMONDS][0] !=
		card.STANDARD_52_RANK_ACE {
		t.Fatal("Expected first diamond to be ace, got",
			ranksBySuit[card.STANDARD_52_SUIT_DIAMONDS][0])
	}
}

func TestIsStraight(t *testing.T) {
	ranks := []int{2, 6, 3, 8, 6}
	ok, _ := IsStraight(ranks)
	if ok {
		t.Fatal("Detected as straight but isn't")
	}
	ranks = []int{2, 6, 3, 4, 5}
	ok, highCard := IsStraight(ranks)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 6 {
		t.Fatal("Expected high card of 6, got", highCard)
	}
	ranks = []int{2, 6, 3, 4, 5, 4}
	ok, highCard = IsStraight(ranks)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 6 {
		t.Fatal("Expected high card of 6, got", highCard)
	}
	// Ace as low card
	ranks = []int{2, 1, 3, 5, 4}
	ok, highCard = IsStraight(ranks)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 5 {
		t.Fatal("Expected high card of 5, got", highCard)
	}
	// Ace as high card
	ranks = []int{11, 10, 13, 12, 1}
	ok, highCard = IsStraight(ranks)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if highCard != 14 {
		t.Fatal("Expected high card of 14, got", highCard)
	}
}

func TestStraightFlush(t *testing.T) {

}
