package modern_art

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mockGame(t *testing.T) *Game {
	players := []string{"Mick", "Steve", "BJ", "Elva"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Fatal(err)
	}
	return game
}

func cloneGame(g *Game) *Game {
	newG := &Game{}
	data, err := g.Encode()
	if err != nil {
		panic(err.Error())
	}
	if err := newG.Decode(data); err != nil {
		panic(err.Error())
	}
	return newG
}

func TestDeck(t *testing.T) {
	Convey("Given a fresh deck", t, func() {
		d := Deck()
		Convey("It should have 70 cards", func() {
			So(len(d), ShouldEqual, 70)
		})
		Convey("It should have 12 Lite Metal cards", func() {
			i := 0
			for _, c := range d.ToSuitRankCards() {
				if c.Suit == SUIT_LITE_METAL {
					i += 1
				}
			}
			So(i, ShouldEqual, 12)
		})
		Convey("It should have 13 Yoko cards", func() {
			i := 0
			for _, c := range d.ToSuitRankCards() {
				if c.Suit == SUIT_YOKO {
					i += 1
				}
			}
			So(i, ShouldEqual, 13)
		})
		Convey("It should have 14 Christine P cards", func() {
			i := 0
			for _, c := range d.ToSuitRankCards() {
				if c.Suit == SUIT_CHRISTINE_P {
					i += 1
				}
			}
			So(i, ShouldEqual, 14)
		})
		Convey("It should have 15 Karl Glitter cards", func() {
			i := 0
			for _, c := range d.ToSuitRankCards() {
				if c.Suit == SUIT_KARL_GLITTER {
					i += 1
				}
			}
			So(i, ShouldEqual, 15)
		})
		Convey("It should have 16 Krypto cards", func() {
			i := 0
			for _, c := range d.ToSuitRankCards() {
				if c.Suit == SUIT_KRYPTO {
					i += 1
				}
			}
			So(i, ShouldEqual, 16)
		})
	})
}

func TestStart(t *testing.T) {
	Convey("Game should start", t, func() {
		mockGame(t)
	})
}
