package poker

import (
	"github.com/beefsack/brdg.me/game/card"
)

const (
	CATEGORY_NONE = iota
	CATEGORY_HIGH_CARD
	CATEGORY_ONE_PAIR
	CATEGORY_TWO_PAIR
	CATEGORY_THREE_OF_A_KIND
	CATEGORY_STRAIGHT
	CATEGORY_FLUSH
	CATEGORY_FULL_HOUSE
	CATEGORY_FOUR_OF_A_KIND
	CATEGORY_STRAIGHT_FLUSH
)

type HandResult struct {
	Category int
	Cards    card.Deck
}

func Result(hand card.Deck) (result HandResult) {
	cardsBySuit := CardsBySuit(hand)
	// Straight flush
	for i := 0; i < 4; i++ {
		ok, cards := IsStraight(cardsBySuit[i])
		if ok && (result.Category < CATEGORY_STRAIGHT_FLUSH ||
			cards[0].(card.SuitRankCard).Rank >
				result.Cards[0].(card.SuitRankCard).Rank) {
			result.Category = CATEGORY_STRAIGHT_FLUSH
			result.Cards = cards
		}
	}
	if result.Category != CATEGORY_NONE {
		return result
	}
	// Four of a kind
	ok, cards := IsFourOfAKind(hand)
	if ok {
		result.Category = CATEGORY_FOUR_OF_A_KIND
		result.Cards = cards
		return result
	}
	// Full house
	ok, cards = IsFullHouse(hand)
	if ok {
		result.Category = CATEGORY_FULL_HOUSE
		result.Cards = cards
		return result
	}
	// Flush
	// Straight
	ok, cards = IsStraight(hand)
	if ok {
		result.Category = CATEGORY_STRAIGHT
		result.Cards = cards
		return result
	}
	// Three of a kind
	ok, cards = IsThreeOfAKind(hand)
	if ok {
		result.Category = CATEGORY_THREE_OF_A_KIND
		result.Cards = cards
		return result
	}
	// Two pair
	ok, cards = IsTwoPair(hand)
	if ok {
		result.Category = CATEGORY_TWO_PAIR
		result.Cards = cards
		return result
	}
	// One pair
	ok, cards = IsOnePair(hand)
	if ok {
		result.Category = CATEGORY_ONE_PAIR
		result.Cards = cards
		return result
	}
	// High card
	result.Category = CATEGORY_HIGH_CARD
	cards, _ = FindHighestRank(hand, 5)
	result.Cards = cards
	return result
}

func IsStraight(hand card.Deck) (ok bool, cards card.Deck) {
	if len(hand) < 5 {
		// Impossible to have a straight with less than five cards
		return false, card.Deck{}
	}
	byRank := CardsByRank(hand)
	for i := card.STANDARD_52_RANK_ACE_HIGH; i >= 2; i-- {
		if len(byRank[i]) > 0 {
			cards = cards.Push(byRank[i][0])
			if len(cards) == 5 {
				ok = true
				break
			}
		} else {
			cards = card.Deck{}
		}
	}
	if len(cards) == 4 && len(byRank[card.STANDARD_52_RANK_ACE_HIGH]) > 0 {
		// Ace also counts as low
		ok = true
		cards = cards.Push(byRank[card.STANDARD_52_RANK_ACE_HIGH][0])
	}
	return
}

func IsFourOfAKind(hand card.Deck) (ok bool, cards card.Deck) {
	ok, cards, remaining := FindMultiple(hand, 4)
	if ok {
		kicker, _ := FindHighestRank(remaining, 1)
		cards = cards.PushMany(kicker)
	}
	return
}

func IsFullHouse(hand card.Deck) (ok bool, cards card.Deck) {
	ok, cards, remaining := FindMultiple(hand, 3)
	if ok {
		ok, pair, _ := FindMultiple(remaining, 2)
		if ok {
			cards = cards.PushMany(pair)
		}
	}
	return
}

func IsThreeOfAKind(hand card.Deck) (ok bool, cards card.Deck) {
	ok, cards, remaining := FindMultiple(hand, 3)
	if ok {
		kickers, _ := FindHighestRank(remaining, 2)
		cards = cards.PushMany(kickers)
	}
	return
}

func IsTwoPair(hand card.Deck) (ok bool, cards card.Deck) {
	ok, cards, remaining := FindMultiple(hand, 2)
	if ok {
		ok, pair, remaining := FindMultiple(remaining, 2)
		if ok {
			cards = cards.PushMany(pair)
			kicker, _ := FindHighestRank(remaining, 1)
			cards = cards.PushMany(kicker)
		}
	}
	return
}

func IsOnePair(hand card.Deck) (ok bool, cards card.Deck) {
	ok, cards, remaining := FindMultiple(hand, 2)
	if ok {
		kickers, _ := FindHighestRank(remaining, 3)
		cards = cards.PushMany(kickers)
	}
	return
}

// Finds a multiple of a rank of card
func FindMultiple(hand card.Deck, n int) (ok bool, cards card.Deck,
	remaining card.Deck) {
	remaining = hand
	byRank := CardsByRank(remaining)
	for i := len(byRank) - 1; i >= 0; i-- {
		if len(byRank[i]) >= n {
			ok = true
			cards = byRank[i][:n]
			for _, c := range cards {
				remaining, _ = remaining.Remove(c, 1)
			}
			break
		}
	}
	return
}

// Pick the highest ranking n cards given the hand
func FindHighestRank(hand card.Deck, n int) (highest card.Deck,
	remaining card.Deck) {
	remaining = hand
	byRank := CardsByRank(remaining)
	for i := len(byRank) - 1; i >= 0; i-- {
		take := n - len(highest)
		if len(byRank[i]) < take {
			take = len(byRank[i])
		}
		highest = highest.PushMany(byRank[i][:take])
		if len(highest) == n {
			break
		}
	}
	return
}

// Breaks down a deck to ranks by suit, sorted by rank ascending
func CardsBySuit(hand card.Deck) map[int]card.Deck {
	ranksBySuit := map[int]card.Deck{}
	// Initialise
	for i := card.STANDARD_52_SUIT_CLUBS; i < card.STANDARD_52_SUIT_SPADES; i++ {
		ranksBySuit[i] = card.Deck{}
	}
	// Categorise
	for _, c := range hand {
		s := c.(card.SuitRankCard).Suit
		ranksBySuit[s] = ranksBySuit[s].Push(c)
	}
	// Sort
	for i := card.STANDARD_52_SUIT_CLUBS; i < card.STANDARD_52_SUIT_SPADES; i++ {
		ranksBySuit[i] = ranksBySuit[i].Sort()
	}
	return ranksBySuit
}

func CardsByRank(hand card.Deck) map[int]card.Deck {
	suitsByRank := map[int]card.Deck{}
	// Initialise
	for i := card.STANDARD_52_RANK_2; i < card.STANDARD_52_RANK_ACE_HIGH; i++ {
		suitsByRank[i] = card.Deck{}
	}
	// Categorise
	for _, c := range hand {
		r := c.(card.SuitRankCard).Rank
		suitsByRank[r] = suitsByRank[r].Push(c)
	}
	return suitsByRank
}
