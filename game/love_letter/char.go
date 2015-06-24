package love_letter

const (
	Princess = 8 - iota
	Countess
	King
	Prince
	Handmaid
	Baron
	Priest
	Guard
)

type Char interface {
	Name() string
	Number() int
	Text() string
	Play(g *Game, player int, args ...string) error
}

type DiscardHandler interface {
	HandleDiscard()
}

var Chars = map[int]Char{
	Princess: CharPrincess{},
	Countess: CharCountess{},
	King:     CharKing{},
	Prince:   CharPrince{},
	Handmaid: CharHandmaid{},
	Baron:    CharBaron{},
	Priest:   CharPriest{},
	Guard:    CharGuard{},
}

var Deck = []int{
	Guard,
	Guard,
	Guard,
	Guard,
	Guard,
	Priest,
	Priest,
	Baron,
	Baron,
	Handmaid,
	Handmaid,
	Prince,
	Prince,
	King,
	Countess,
	Princess,
}
