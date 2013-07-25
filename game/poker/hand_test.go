package poker

import (
	"github.com/beefsack/brdg.me/game/card"
	"testing"
)

func buildHandByRanks(ranks []int) card.Deck {
	d := card.Deck{}
	for _, r := range ranks {
		d = d.Push(card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: r,
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
	ok, _ := IsStraight(hand)
	if ok {
		t.Fatal("Detected as straight but isn't")
	}
	hand = buildHandByRanks([]int{2, 6, 3, 4, 5})
	ok, cards := IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if cards[0].(card.SuitRankCard).Rank != 6 {
		t.Fatal("Expected high card of 6, got", cards[0].(card.SuitRankCard).Rank)
	}
	hand = buildHandByRanks([]int{2, 6, 3, 4, 5, 4})
	ok, cards = IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if cards[0].(card.SuitRankCard).Rank != 6 {
		t.Fatal("Expected high card of 6, got", cards[0].(card.SuitRankCard).Rank)
	}
	// Ace as low card
	hand = buildHandByRanks([]int{2, 14, 3, 5, 4})
	ok, cards = IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if cards[0].(card.SuitRankCard).Rank != 5 {
		t.Fatal("Expected high card of 5, got", cards[0].(card.SuitRankCard).Rank)
	}
	// Ace as high card
	hand = buildHandByRanks([]int{11, 10, 13, 12, 14})
	ok, cards = IsStraight(hand)
	if !ok {
		t.Fatal("Didn't detect as straight")
	}
	if cards[0].(card.SuitRankCard).Rank != 14 {
		t.Fatal("Expected high card of 14, got", cards[0].(card.SuitRankCard).Rank)
	}
}

func TestStraightFlush(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_7,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_5,
		},
	})
	if handResult.Category != CATEGORY_STRAIGHT_FLUSH {
		t.Fatal("Expected straight flush, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_7 {
		t.Fatal("Expected 7 high, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
}

func TestFourOfAKind(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_HEARTS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_5,
		},
	})
	if handResult.Category != CATEGORY_FOUR_OF_A_KIND {
		t.Fatal("Expected four of a kind, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_3 {
		t.Fatal("Expected first rank of 3, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[4].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_6 {
		t.Fatal("Expected fourth rank of 6, got:",
			handResult.Cards[4].(card.SuitRankCard).Rank)
	}
}

func TestFullHouse(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_HEARTS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_5,
		},
	})
	if handResult.Category != CATEGORY_FULL_HOUSE {
		t.Fatal("Expected full house, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_3 {
		t.Fatal("Expected first rank of 3, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[3].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_6 {
		t.Fatal("Expected second rank of 6, got:",
			handResult.Cards[3].(card.SuitRankCard).Rank)
	}
}

func TestStraight(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_2,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_ACE_HIGH,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_5,
		},
	})
	if handResult.Category != CATEGORY_STRAIGHT {
		t.Fatal("Expected straight, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_5 {
		t.Fatal("Expected 5 high, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
}

func TestThreeOfAKind(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_2,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_5,
		},
	})
	if handResult.Category != CATEGORY_THREE_OF_A_KIND {
		t.Fatal("Expected three of a kind, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank != card.STANDARD_52_RANK_3 {
		t.Fatal("Expected first card to be 3, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
}

func TestTwoPair(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_2,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_KING,
		},
	})
	if handResult.Category != CATEGORY_TWO_PAIR {
		t.Fatal("Expected two pair, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_KING {
		t.Fatal("Expected first card to be king, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[2].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_3 {
		t.Fatal("Expected third card to be 3, got:",
			handResult.Cards[2].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[4].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_6 {
		t.Fatal("Expected fifth card to be 6, got:",
			handResult.Cards[4].(card.SuitRankCard).Rank)
	}
}

func TestOnePair(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_2,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_9,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_KING,
		},
	})
	if handResult.Category != CATEGORY_ONE_PAIR {
		t.Fatal("Expected one pair, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_KING {
		t.Fatal("Expected first card to be king, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[2].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_9 {
		t.Fatal("Expected third card to be 9, got:",
			handResult.Cards[2].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[3].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_6 {
		t.Fatal("Expected fourth card to be 6, got:",
			handResult.Cards[3].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[4].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_4 {
		t.Fatal("Expected fifth card to be 4, got:",
			handResult.Cards[4].(card.SuitRankCard).Rank)
	}
}

func TestHighCard(t *testing.T) {
	handResult := Result(card.Deck{
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_2,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_3,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_KING,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_SPADES,
			Rank: card.STANDARD_52_RANK_6,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_4,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_CLUBS,
			Rank: card.STANDARD_52_RANK_9,
		},
		card.SuitRankCard{
			Suit: card.STANDARD_52_SUIT_DIAMONDS,
			Rank: card.STANDARD_52_RANK_QUEEN,
		},
	})
	if handResult.Category != CATEGORY_HIGH_CARD {
		t.Fatal("Expected high card, got:", handResult.Category)
	}
	if handResult.Cards[0].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_KING {
		t.Fatal("Expected first card to be king, got:",
			handResult.Cards[0].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[1].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_QUEEN {
		t.Fatal("Expected second card to be queen, got:",
			handResult.Cards[1].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[2].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_9 {
		t.Fatal("Expected third card to be 9, got:",
			handResult.Cards[2].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[3].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_6 {
		t.Fatal("Expected fourth card to be 6, got:",
			handResult.Cards[3].(card.SuitRankCard).Rank)
	}
	if handResult.Cards[4].(card.SuitRankCard).Rank !=
		card.STANDARD_52_RANK_4 {
		t.Fatal("Expected fifth card to be 4, got:",
			handResult.Cards[4].(card.SuitRankCard).Rank)
	}
}
