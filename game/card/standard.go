package card

const (
	STANDARD_52_SUIT_CLUBS = iota
	STANDARD_52_SUIT_DIAMONDS
	STANDARD_52_SUIT_HEARTS
	STANDARD_52_SUIT_SPADES
	STANDARD_52_SUIT_JOKER
)

const (
	_ = iota // Ignore 0
	STANDARD_52_VALUE_ACE
	STANDARD_52_VALUE_2
	STANDARD_52_VALUE_3
	STANDARD_52_VALUE_4
	STANDARD_52_VALUE_5
	STANDARD_52_VALUE_6
	STANDARD_52_VALUE_7
	STANDARD_52_VALUE_8
	STANDARD_52_VALUE_9
	STANDARD_52_VALUE_10
	STANDARD_52_VALUE_JACK
	STANDARD_52_VALUE_QUEEN
	STANDARD_52_VALUE_KING
)

func Standard52Deck() Deck {
	d := Deck{}
	for suit := STANDARD_52_SUIT_CLUBS; suit <= STANDARD_52_SUIT_SPADES; suit++ {
		for value := STANDARD_52_VALUE_ACE; value <= STANDARD_52_VALUE_KING; value++ {
			d = append(d, SuitValueCard{
				Suit:  suit,
				Value: value,
			})
		}
	}
	return d
}

func Standard52DeckWithJokers() (d Deck) {
	d = Standard52Deck()
	for i := 0; i < 2; i++ {
		d = append(d, SuitValueCard{
			Suit: STANDARD_52_SUIT_JOKER,
		})
	}
	return d
}
