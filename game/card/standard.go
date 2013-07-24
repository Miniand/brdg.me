package card

import (
	"fmt"
)

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

func (c SuitValueCard) RenderStandard52() string {
	var (
		symbol string
		colour string
		value  string
	)
	switch c.Suit {
	case STANDARD_52_SUIT_CLUBS:
		symbol = "♣"
		colour = "black"
	case STANDARD_52_SUIT_DIAMONDS:
		symbol = "♦"
		colour = "red"
	case STANDARD_52_SUIT_HEARTS:
		symbol = "♥"
		colour = "red"
	case STANDARD_52_SUIT_SPADES:
		symbol = "♠"
		colour = "black"
	}
	switch c.Value {
	case STANDARD_52_VALUE_ACE:
		value = "A"
	case STANDARD_52_VALUE_JACK:
		value = "J"
	case STANDARD_52_VALUE_QUEEN:
		value = "Q"
	case STANDARD_52_VALUE_KING:
		value = "K"
	default:
		value = fmt.Sprintf("%d", c.Value)
	}
	return fmt.Sprintf(`{{c "%s"}}%s%s{{_c}}`, colour, symbol, value)
}

func (c SuitValueCard) RenderStandard52FixedWidth() string {
	output := c.RenderStandard52()
	if c.Value != 10 {
		output += " "
	}
	return output
}

func RenderStandard52Hidden() string {
	return `{{c "gray"}}##{{_c}}`
}

func RenderStandard52HiddenFixedWidth() string {
	return RenderStandard52Hidden() + " "
}
