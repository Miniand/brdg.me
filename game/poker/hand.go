package poker

import (
	"fmt"
	"github.com/beefsack/brdg.me/game/card"
	"sort"
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
	Ranks    []int
}

func Result(hand card.Deck) HandResult {
	handResult := HandResult{}
	// suitsByRank := SuitsByRank(hand)
	ranksBySuit := SuitsByRank(hand)
	// Straight flush
	for i := 0; i < 4; i++ {
		if len(ranksBySuit[i]) >= 5 {
			ok, highCard := IsStraight(ranksBySuit[i])
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

func IsStraight(ranks []int) (bool, int) {
	// Work on a copy so we don't break the original
	ranksCopy := make([]int, len(ranks))
	copy(ranksCopy, ranks)
	ranksCopy = SortRanks(ranksCopy)
	fmt.Println(ranksCopy)
	if len(ranksCopy) > 0 && ranksCopy[0] ==
		card.STANDARD_52_RANK_ACE {
		// Ace can be the high card of a straight too
		fmt.Println(ranksCopy)
		ranksCopy = append([]int{14}, ranksCopy...)
		fmt.Println(ranksCopy)
	}
	highCard := 0
	consecutive := 0
	for _, r := range ranksCopy {
		if highCard == 0 || (r != highCard-consecutive+1 &&
			r != highCard-consecutive) {
			// Reset
			highCard = r
			consecutive = 1
		} else if r == highCard-consecutive {
			// Consecutive card
			consecutive++
			if consecutive == 5 {
				return true, highCard
			}
		}
	}
	return false, 0
}

// Breaks down a deck to ranks by suit, sorted by rank descending
func RanksBySuit(hand card.Deck) map[int][]int {
	ranksBySuit := map[int][]int{}
	// Initialise
	for i := 0; i < 4; i++ {
		ranksBySuit[i] = []int{}
	}
	// Categorise
	for _, c := range hand {
		s := c.(card.SuitRankCard).Suit
		r := c.(card.SuitRankCard).Rank
		ranksBySuit[s] = append(ranksBySuit[s], r)
	}
	// Sort
	for i := 0; i < 4; i++ {
		ranksBySuit[i] = SortRanks(ranksBySuit[i])
	}
	return ranksBySuit
}

func SuitsByRank(hand card.Deck) map[int][]int {
	suitsByRank := map[int][]int{}
	// Initialise
	for i := 0; i < 14; i++ {
		suitsByRank[i] = []int{}
	}
	// Categorise
	for _, c := range hand {
		s := c.(card.SuitRankCard).Suit
		r := c.(card.SuitRankCard).Rank
		suitsByRank[r] = append(suitsByRank[r], s)
	}
	return suitsByRank
}

func SortRanks(ranks []int) []int {
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))
	if len(ranks) > 0 && ranks[0] != card.STANDARD_52_RANK_ACE {
		// Bring aces to the front as they are more valuable in poker
		for ranks[len(ranks)-1] == card.STANDARD_52_RANK_ACE {
			ranks = append([]int{card.STANDARD_52_RANK_KING + 1},
				ranks[:len(ranks)-1]...)
		}
	}
	return ranks
}
