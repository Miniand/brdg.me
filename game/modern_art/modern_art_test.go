package modern_art

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

const (
	MICK = iota
	STEVE
	BJ
	ELVA
)

var playerNames = map[int]string{
	MICK:  "Mick",
	STEVE: "Steve",
	BJ:    "BJ",
	ELVA:  "Elva",
}

func mockGame(t *testing.T) *Game {
	players := []string{
		playerNames[MICK],
		playerNames[STEVE],
		playerNames[BJ],
		playerNames[ELVA],
	}
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
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("It should have given each player 9 cards for 4 players", func() {
			So(len(g.PlayerHands[0]), ShouldEqual, 9)
			So(len(g.PlayerHands[1]), ShouldEqual, 9)
			So(len(g.PlayerHands[2]), ShouldEqual, 9)
			So(len(g.PlayerHands[3]), ShouldEqual, 9)
		})
		Convey("It should have left 34 cards in the deck", func() {
			So(len(g.Deck), ShouldEqual, 34)
		})
		Convey("It should have given $100 to each player", func() {
			So(g.PlayerMoney[0], ShouldEqual, 100)
			So(g.PlayerMoney[1], ShouldEqual, 100)
			So(g.PlayerMoney[2], ShouldEqual, 100)
			So(g.PlayerMoney[3], ShouldEqual, 100)
		})
	})
}

func TestOpenAuction(t *testing.T) {
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("Given BJ has a Lite Metal Open Auction card", func() {
			g := cloneGame(g)
			g.CurrentPlayer = BJ
			g.PlayerHands[BJ] = g.PlayerHands[BJ].Push(card.SuitRankCard{
				SUIT_LITE_METAL, RANK_OPEN})
			Convey("Given BJ plays the Lite Metal Open Auction card", func() {
				g := cloneGame(g)
				_, err := command.CallInCommands(playerNames[BJ], g,
					"play lmop", g.Commands(playerNames[BJ]))
				So(err, ShouldBeNil)
				So(g.State, ShouldEqual, STATE_AUCTION)
				So(len(g.CurrentlyAuctioning), ShouldEqual, 1)
				Convey("Given Steve bids", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[STEVE], g,
						"bid 10", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					So(g.State, ShouldEqual, STATE_AUCTION)
					Convey("Given the other players all pass", func() {
						g := cloneGame(g)
						_, err := command.CallInCommands(playerNames[MICK], g,
							"pass", g.Commands(playerNames[MICK]))
						So(err, ShouldBeNil)
						_, err = command.CallInCommands(playerNames[BJ], g,
							"pass", g.Commands(playerNames[BJ]))
						So(err, ShouldBeNil)
						_, err = command.CallInCommands(playerNames[ELVA], g,
							"pass", g.Commands(playerNames[ELVA]))
						So(err, ShouldBeNil)
						Convey("It should give the card to Steve and go to the next player", func() {
							So(g.State, ShouldEqual, STATE_PLAY_CARD)
							So(g.CurrentPlayer, ShouldEqual, ELVA)
							So(len(g.PlayerPurchases[STEVE]), ShouldEqual, 1)
							So(g.PlayerMoney[STEVE], ShouldEqual, 90)
							So(g.PlayerMoney[BJ], ShouldEqual, 110)
						})
					})
				})
				Convey("Given nobody bids", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[MICK], g,
						"pass", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[STEVE], g,
						"pass", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[ELVA], g,
						"pass", g.Commands(playerNames[ELVA]))
					So(err, ShouldBeNil)
					Convey("It should give BJ the card for nothing", func() {
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, ELVA)
						So(len(g.PlayerPurchases[BJ]), ShouldEqual, 1)
						So(g.PlayerMoney[BJ], ShouldEqual, 100)
					})
				})
			})
		})
	})
}

func TestFixedPriceAuction(t *testing.T) {
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("Given Elva has a Christine P Fixed Price Auction card", func() {
			g := cloneGame(g)
			g.CurrentPlayer = ELVA
			g.PlayerHands[ELVA] = g.PlayerHands[ELVA].Push(card.SuitRankCard{
				SUIT_CHRISTINE_P, RANK_FIXED_PRICE})
			Convey("Given Elva plays the Christine P Fixed Price Auction card and sets the price at 15", func() {
				g := cloneGame(g)
				_, err := command.CallInCommands(playerNames[ELVA], g,
					"play cpfp", g.Commands(playerNames[ELVA]))
				So(err, ShouldBeNil)
				So(g.State, ShouldEqual, STATE_AUCTION)
				So(len(g.CurrentlyAuctioning), ShouldEqual, 1)
				_, err = command.CallInCommands(playerNames[ELVA], g,
					"price 15", g.Commands(playerNames[ELVA]))
				So(err, ShouldBeNil)
				Convey("Given Mick passes and Steve buys", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[MICK], g,
						"pass", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					So(g.State, ShouldEqual, STATE_AUCTION)
					_, err = command.CallInCommands(playerNames[STEVE], g,
						"buy", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					Convey("Steve should receive the card for the given price", func() {
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, MICK)
						So(len(g.PlayerPurchases[STEVE]), ShouldEqual, 1)
						So(g.PlayerMoney[STEVE], ShouldEqual, 85)
						So(g.PlayerMoney[ELVA], ShouldEqual, 115)
					})
				})
				Convey("Given nobody bids", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[MICK], g,
						"pass", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[STEVE], g,
						"pass", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[BJ], g,
						"pass", g.Commands(playerNames[BJ]))
					So(err, ShouldBeNil)
					Convey("It should give the card to Elva for the given price", func() {
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, MICK)
						So(len(g.PlayerPurchases[ELVA]), ShouldEqual, 1)
						So(g.PlayerMoney[ELVA], ShouldEqual, 85)
					})
				})
			})
		})
	})
}

func TestSealedAuction(t *testing.T) {
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("Given Elva has a Krypto Sealed Auction card", func() {
			g := cloneGame(g)
			g.CurrentPlayer = ELVA
			g.PlayerHands[ELVA] = g.PlayerHands[ELVA].Push(card.SuitRankCard{
				SUIT_KRYPTO, RANK_SEALED})
			Convey("Given Elva plays the Krypto Sealed Auction card", func() {
				g := cloneGame(g)
				_, err := command.CallInCommands(playerNames[ELVA], g,
					"play krsl", g.Commands(playerNames[ELVA]))
				So(err, ShouldBeNil)
				So(g.State, ShouldEqual, STATE_AUCTION)
				So(len(g.CurrentlyAuctioning), ShouldEqual, 1)
				Convey("Given everyone bids different amounts", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[MICK], g,
						"bid 4", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[STEVE], g,
						"bid 5", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[BJ], g,
						"bid 3", g.Commands(playerNames[BJ]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[ELVA], g,
						"bid 1", g.Commands(playerNames[ELVA]))
					So(err, ShouldBeNil)
					Convey("Steve should receive the card for the given price", func() {
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, MICK)
						So(len(g.PlayerPurchases[STEVE]), ShouldEqual, 1)
						So(g.PlayerMoney[STEVE], ShouldEqual, 95)
						So(g.PlayerMoney[ELVA], ShouldEqual, 105)
					})
				})
				Convey("Given nobody bids", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[MICK], g,
						"pass", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[STEVE], g,
						"pass", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[ELVA], g,
						"pass", g.Commands(playerNames[ELVA]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[BJ], g,
						"pass", g.Commands(playerNames[BJ]))
					So(err, ShouldBeNil)
					Convey("It should give the card to Elva for free", func() {
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, MICK)
						So(len(g.PlayerPurchases[ELVA]), ShouldEqual, 1)
						So(g.PlayerMoney[ELVA], ShouldEqual, 100)
					})
				})
			})
		})
	})
}

func TestDoubleAuction(t *testing.T) {
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("Given Elva has a Karl Glitter Double Auction card and Steve has a Karl Glitter Sealed Auction card", func() {
			g := cloneGame(g)
			g.CurrentPlayer = ELVA
			g.PlayerHands[ELVA] = g.PlayerHands[ELVA].Push(card.SuitRankCard{
				SUIT_KARL_GLITTER, RANK_DOUBLE}).Push(card.SuitRankCard{
				SUIT_KARL_GLITTER, RANK_SEALED})
			g.PlayerHands[STEVE] = g.PlayerHands[STEVE].Push(card.SuitRankCard{
				SUIT_KARL_GLITTER, RANK_SEALED})
			Convey("Given Elva plays the Karl Glitter Double Auction card", func() {
				g := cloneGame(g)
				_, err := command.CallInCommands(playerNames[ELVA], g,
					"play kgdb", g.Commands(playerNames[ELVA]))
				So(err, ShouldBeNil)
				So(g.State, ShouldEqual, STATE_AUCTION)
				So(len(g.CurrentlyAuctioning), ShouldEqual, 1)
				Convey("Given Elva passes, Mick passes and Steve plays his KG Sealed", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[ELVA], g,
						"pass", g.Commands(playerNames[ELVA]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[MICK], g,
						"pass", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[STEVE], g,
						"add kgsl", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					Convey("It should start a new sealed auction with Steve as the auctioneer", func() {
						g := cloneGame(g)
						So(g.CurrentPlayer, ShouldEqual, STEVE)
						So(len(g.CurrentlyAuctioning), ShouldEqual, 2)
						So(g.IsAuction(), ShouldBeTrue)
						So(g.AuctionType(), ShouldEqual, RANK_SEALED)
						Convey("Given everyone bids different amounts", func() {
							g := cloneGame(g)
							_, err := command.CallInCommands(playerNames[MICK], g,
								"bid 8", g.Commands(playerNames[MICK]))
							So(err, ShouldBeNil)
							_, err = command.CallInCommands(playerNames[STEVE], g,
								"bid 5", g.Commands(playerNames[STEVE]))
							So(err, ShouldBeNil)
							_, err = command.CallInCommands(playerNames[BJ], g,
								"bid 3", g.Commands(playerNames[BJ]))
							So(err, ShouldBeNil)
							_, err = command.CallInCommands(playerNames[ELVA], g,
								"bid 1", g.Commands(playerNames[ELVA]))
							So(err, ShouldBeNil)
							Convey("Mick should receive both the cards for the given price", func() {
								So(g.State, ShouldEqual, STATE_PLAY_CARD)
								So(g.CurrentPlayer, ShouldEqual, BJ)
								So(len(g.PlayerPurchases[MICK]), ShouldEqual, 2)
								So(g.PlayerMoney[MICK], ShouldEqual, 92)
								So(g.PlayerMoney[STEVE], ShouldEqual, 108)
								So(g.PlayerMoney[ELVA], ShouldEqual, 100)
							})
						})
					})
				})
			})
		})
	})
}

func TestDoubleAuctionEndsRound(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.CurrentPlayer = helper.Mick
	g.PlayerPurchases = map[int]card.Deck{
		helper.Mick: card.Deck{
			card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
		},
		helper.Steve: card.Deck{
			card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
		},
		helper.BJ: card.Deck{
			card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
		},
	}
	g.PlayerHands[helper.Mick] = card.Deck{
		card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
		card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
		card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
		card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE},
	}
	g.PlayerHands[helper.Steve] = card.Deck{
		card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
		card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
		card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
		card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
	}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play lmdb"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "pass"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "add lmop"))
	assert.Equal(t, 1, g.Round)
	assert.Equal(t, helper.BJ, g.CurrentPlayer)
}

func TestOnceAroundAuction(t *testing.T) {
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("Given Mick has a Yoko Once Around Auction card", func() {
			g := cloneGame(g)
			g.CurrentPlayer = MICK
			g.PlayerHands[MICK] = g.PlayerHands[MICK].Push(card.SuitRankCard{
				SUIT_YOKO, RANK_ONCE_AROUND})
			Convey("Given Mick plays the Yoko Once Around Auction card", func() {
				g := cloneGame(g)
				_, err := command.CallInCommands(playerNames[MICK], g,
					"play yooa", g.Commands(playerNames[MICK]))
				So(err, ShouldBeNil)
				So(g.State, ShouldEqual, STATE_AUCTION)
				So(len(g.CurrentlyAuctioning), ShouldEqual, 1)
				Convey("Given some bids", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[STEVE], g,
						"pass", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[BJ], g,
						"bid 5", g.Commands(playerNames[BJ]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[ELVA], g,
						"bid 7", g.Commands(playerNames[ELVA]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[MICK], g,
						"pass", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					Convey("It should give the card to Elva", func() {
						g := cloneGame(g)
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, STEVE)
						So(len(g.PlayerPurchases[ELVA]), ShouldEqual, 1)
						So(g.PlayerMoney[MICK], ShouldEqual, 107)
						So(g.PlayerMoney[ELVA], ShouldEqual, 93)
					})
				})
				Convey("Given everyone passes", func() {
					g := cloneGame(g)
					_, err := command.CallInCommands(playerNames[STEVE], g,
						"pass", g.Commands(playerNames[STEVE]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[BJ], g,
						"pass", g.Commands(playerNames[BJ]))
					So(err, ShouldBeNil)
					_, err = command.CallInCommands(playerNames[ELVA], g,
						"pass", g.Commands(playerNames[ELVA]))
					So(err, ShouldBeNil)
					Convey("It should give the card to Mick for free", func() {
						g := cloneGame(g)
						So(g.State, ShouldEqual, STATE_PLAY_CARD)
						So(g.CurrentPlayer, ShouldEqual, STEVE)
						So(len(g.PlayerPurchases[MICK]), ShouldEqual, 1)
						So(g.PlayerMoney[MICK], ShouldEqual, 100)
					})
				})
			})
		})
	})
}

func TestEndOfRound(t *testing.T) {
	Convey("Given a new game", t, func() {
		g := mockGame(t)
		Convey("Given there are already 3 Lite Metal on the board", func() {
			g := cloneGame(g)
			g.PlayerPurchases[MICK] = card.Deck{
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
			}
			g.PlayerPurchases[STEVE] = card.Deck{
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
			}
			Convey("Given Mick plays a Lite Metal Double Auction", func() {
				g := cloneGame(g)
				g.PlayerHands[MICK] = g.PlayerHands[MICK].Push(
					card.SuitRankCard{SUIT_LITE_METAL, RANK_DOUBLE})
				_, err := command.CallInCommands(playerNames[MICK], g,
					"play lmdb", g.Commands(playerNames[MICK]))
				So(err, ShouldBeNil)
				Convey("It should be the same round", func() {
					g := cloneGame(g)
					So(g.Round, ShouldEqual, 0)
				})
				Convey("Given Mick adds another Lite Metal", func() {
					g := cloneGame(g)
					g.PlayerHands[MICK] = g.PlayerHands[MICK].Push(
						card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN})
					_, err := command.CallInCommands(playerNames[MICK], g,
						"add lmop", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					Convey("It should be the next round and values should be added to artists", func() {
						g := cloneGame(g)
						So(g.Round, ShouldEqual, 1)
						So(g.SuitValue(SUIT_LITE_METAL), ShouldEqual, 30)
						So(g.PlayerMoney[MICK], ShouldEqual, 160)
						So(g.PlayerMoney[STEVE], ShouldEqual, 130)
						So(g.PlayerMoney[BJ], ShouldEqual, 100)
						So(g.PlayerMoney[ELVA], ShouldEqual, 100)
					})
				})
			})
		})
		Convey("Given there are already 4 Lite Metal on the board", func() {
			g := cloneGame(g)
			g.PlayerPurchases[MICK] = card.Deck{
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
			}
			g.PlayerPurchases[STEVE] = card.Deck{
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
			}
			g.PlayerPurchases[BJ] = card.Deck{
				card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN},
			}
			Convey("Given Mick adds another Lite Metal", func() {
				g := cloneGame(g)
				g.PlayerHands[MICK] = g.PlayerHands[MICK].Push(
					card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN})
				_, err := command.CallInCommands(playerNames[MICK], g,
					"play lmop", g.Commands(playerNames[MICK]))
				So(err, ShouldBeNil)
				Convey("It should be the next round and values should be added to artists", func() {
					g := cloneGame(g)
					So(g.Round, ShouldEqual, 1)
					So(g.SuitValue(SUIT_LITE_METAL), ShouldEqual, 30)
					So(g.PlayerMoney[MICK], ShouldEqual, 160)
					So(g.PlayerMoney[STEVE], ShouldEqual, 130)
					So(g.PlayerMoney[BJ], ShouldEqual, 130)
					So(g.PlayerMoney[ELVA], ShouldEqual, 100)
				})
			})
			Convey("Given it is the final round", func() {
				g := cloneGame(g)
				g.Round = 3
				Convey("Given Mick adds another Lite Metal", func() {
					g := cloneGame(g)
					g.PlayerHands[MICK] = g.PlayerHands[MICK].Push(
						card.SuitRankCard{SUIT_LITE_METAL, RANK_OPEN})
					_, err := command.CallInCommands(playerNames[MICK], g,
						"play lmop", g.Commands(playerNames[MICK]))
					So(err, ShouldBeNil)
					Convey("It should be the end of the game and values should be added to artists", func() {
						g := cloneGame(g)
						So(g.IsFinished(), ShouldBeTrue)
						So(g.SuitValue(SUIT_LITE_METAL), ShouldEqual, 30)
						So(g.PlayerMoney[MICK], ShouldEqual, 160)
						So(g.PlayerMoney[STEVE], ShouldEqual, 130)
						So(g.PlayerMoney[BJ], ShouldEqual, 130)
						So(g.PlayerMoney[ELVA], ShouldEqual, 100)
						winners := g.Winners()
						So(len(winners), ShouldEqual, 1)
						So(winners[0], ShouldEqual, playerNames[MICK])
					})
				})
			})
		})
	})
}
