package poker

import (
	"github.com/Miniand/brdg.me/game/card"
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
	Name     string
}

func (hr HandResult) HandScore() []int {
	score := []int{hr.Category}
	for _, c := range hr.Cards {
		score = append(score, c.(card.SuitRankCard).Rank)
	}
	return score
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
		result.Name = "straight flush"
		return result
	}
	// Four of a kind
	ok, cards := IsFourOfAKind(hand)
	if ok {
		result.Category = CATEGORY_FOUR_OF_A_KIND
		result.Cards = cards
		result.Name = "four of a kind"
		return result
	}
	// Full house
	ok, cards = IsFullHouse(hand)
	if ok {
		result.Category = CATEGORY_FULL_HOUSE
		result.Cards = cards
		result.Name = "full house"
		return result
	}
	// Flush
	ok, cards = IsFlush(hand)
	if ok {
		result.Category = CATEGORY_FLUSH
		result.Cards = cards
		result.Name = "flush"
		return result
	}
	// Straight
	ok, cards = IsStraight(hand)
	if ok {
		result.Category = CATEGORY_STRAIGHT
		result.Cards = cards
		result.Name = "straight"
		return result
	}
	// Three of a kind
	ok, cards = IsThreeOfAKind(hand)
	if ok {
		result.Category = CATEGORY_THREE_OF_A_KIND
		result.Cards = cards
		result.Name = "three of a kind"
		return result
	}
	// Two pair
	ok, cards = IsTwoPair(hand)
	if ok {
		result.Category = CATEGORY_TWO_PAIR
		result.Cards = cards
		result.Name = "two pair"
		return result
	}
	// One pair
	ok, cards = IsOnePair(hand)
	if ok {
		result.Category = CATEGORY_ONE_PAIR
		result.Cards = cards
		result.Name = "one pair"
		return result
	}
	// High card
	result.Category = CATEGORY_HIGH_CARD
	cards, _ = FindHighestRank(hand, 5)
	result.Cards = cards
	result.Name = "high card"
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
	var pair card.Deck
	ok, cards, remaining := FindMultiple(hand, 3)
	if ok {
		ok, pair, _ = FindMultiple(remaining, 2)
		if ok {
			cards = cards.PushMany(pair)
		}
	}
	return
}

func IsFlush(hand card.Deck) (ok bool, cards card.Deck) {
	handResults := map[int]HandResult{}
	i := 0
	bySuit := CardsBySuit(hand)
	for suit := card.STANDARD_52_SUIT_CLUBS; suit <=
		card.STANDARD_52_SUIT_SPADES; suit++ {
		if len(bySuit[suit]) >= 5 {
			flush, _ := FindHighestRank(bySuit[suit], 5)
			handResults[i] = HandResult{
				Category: CATEGORY_FLUSH,
				Cards:    flush,
			}
			i++
		}
	}
	if len(handResults) > 0 {
		winningHand := WinningHandResult(handResults)
		if len(winningHand) > 0 {
			ok = true
			cards = handResults[winningHand[0]].Cards
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
	var pair card.Deck
	ok, cards, remaining := FindMultiple(hand, 2)
	if ok {
		ok, pair, remaining = FindMultiple(remaining, 2)
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
	for i := card.STANDARD_52_RANK_ACE_HIGH; i >= 0; i-- {
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
	for i := card.STANDARD_52_RANK_ACE_HIGH; i >= 0; i-- {
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

func WinningHandResult(handResults map[int]HandResult) []int {
	handScores := map[int][]int{}
	nextPass := []int{}
	// Get the scores
	for id, hr := range handResults {
		if hr.Category != CATEGORY_NONE {
			handScores[id] = hr.HandScore()
			nextPass = append(nextPass, id)
		}
	}
	// Loop until we have a winner
	valIndex := 0
	for len(nextPass) > 1 {
		leaders := []int{}
		highest := 0
		for _, handIndex := range nextPass {
			if len(handScores[handIndex]) <= valIndex {
				// Run out of cards, call it a tie
				return nextPass
			}
			if handScores[handIndex][valIndex] > highest {
				leaders = []int{}
				highest = handScores[handIndex][valIndex]
			}
			if handScores[handIndex][valIndex] == highest {
				leaders = append(leaders, handIndex)
			}
		}
		highest = 0
		valIndex++
		nextPass = leaders
	}
	return nextPass
}
