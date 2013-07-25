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
	Ranks    []int
}

func Result(hand card.Deck) HandResult {
	handResult := HandResult{}
	// cardsByRank := CardsByRank(hand)
	cardsBySuit := CardsBySuit(hand)
	// Straight flush
	for i := 0; i < 4; i++ {
		if len(cardsBySuit[i]) >= 5 {
			ok, highCard, _ := IsStraight(cardsBySuit[i])
			if ok && (handResult.Category < CATEGORY_STRAIGHT_FLUSH ||
				highCard > handResult.Ranks[0]) {
				handResult.Category = CATEGORY_STRAIGHT_FLUSH
				handResult.Ranks = []int{highCard}
			}
		}
	}
	// Four of a kind
	// Full house
	// Flush
	// Straight
	// Three of a kind
	// Two pair
	// One pair
	// High card
	return handResult
}

func IsStraight(cards card.Deck) (bool, int, card.Deck) {
	cards = cards.Sort()
	highCard := 0
	consecutive := card.Deck{}
	lastRank := 0
	hasAce := false
	for i := len(cards) - 1; i >= 0; i-- {
		c := cards[i].(card.SuitRankCard)
		if c.Rank == card.STANDARD_52_RANK_ACE {
			hasAce = true
		}
		if highCard == 0 || (c.RankValue() != lastRank-1 && c.Rank != lastRank) {
			// Reset
			highCard = c.RankValue()
			consecutive = card.Deck{c}
		} else if c.RankValue() == lastRank-1 {
			// Consecutive card
			consecutive = consecutive.Unshift(c)
			if len(consecutive) == 5 {
				return true, highCard, consecutive
			}
		}
		lastRank = c.RankValue()
	}
	// Special case if they have the ace and are on a 5 high straight
	if len(consecutive) == 4 && lastRank == card.STANDARD_52_RANK_2 && hasAce {
		return true, highCard, consecutive.Unshift(cards[0])
	}
	return false, 0, consecutive
}

// Breaks down a deck to ranks by suit, sorted by rank descending
func CardsBySuit(hand card.Deck) map[int]card.Deck {
	ranksBySuit := map[int]card.Deck{}
	// Initialise
	for i := 0; i < 4; i++ {
		ranksBySuit[i] = card.Deck{}
	}
	// Categorise
	for _, c := range hand {
		s := c.(card.SuitRankCard).Suit
		ranksBySuit[s] = ranksBySuit[s].Push(c)
	}
	// Sort
	for i := 0; i < 4; i++ {
		ranksBySuit[i] = ranksBySuit[i].Sort()
	}
	return ranksBySuit
}

func CardsByRank(hand card.Deck) map[int]card.Deck {
	suitsByRank := map[int]card.Deck{}
	// Initialise
	for i := 0; i < 14; i++ {
		suitsByRank[i] = card.Deck{}
	}
	// Categorise
	for _, c := range hand {
		r := c.(card.SuitRankCard).Rank
		suitsByRank[r] = suitsByRank[r].Push(c)
	}
	return suitsByRank
}
