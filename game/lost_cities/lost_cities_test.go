package lost_cities

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	. "github.com/smartystreets/goconvey/convey"
)

// Build a game by hand for testing purposes.  Each player has a full hand, half
// of the discard stacks have cards, and there are two cards in the draw pile.
func mockGame(t *testing.T) *Game {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Fatal(err)
	}
	// Mick is the first player
	game.CurrentlyMoving = 0
	// Set Mick's hand
	game.Board.PlayerHands[0] = card.Deck{
		card.SuitRankCard{SUIT_BLUE, 6},
		card.SuitRankCard{SUIT_BLUE, 8},
		card.SuitRankCard{SUIT_RED, 4},
		card.SuitRankCard{SUIT_RED, 5},
		card.SuitRankCard{SUIT_YELLOW, 0},
		card.SuitRankCard{SUIT_YELLOW, 3},
		card.SuitRankCard{SUIT_GREEN, 2},
		card.SuitRankCard{SUIT_WHITE, 10},
	}
	// Set Steve's hand
	game.Board.PlayerHands[1] = card.Deck{
		card.SuitRankCard{SUIT_BLUE, 7},
		card.SuitRankCard{SUIT_BLUE, 9},
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 10},
		card.SuitRankCard{SUIT_YELLOW, 4},
		card.SuitRankCard{SUIT_YELLOW, 7},
		card.SuitRankCard{SUIT_GREEN, 4},
		card.SuitRankCard{SUIT_WHITE, 8},
	}
	// Just set the draw pile to have a couple of cards so we can finish the
	// round quickly for testing.
	game.Board.DrawPile = card.Deck{
		card.SuitRankCard{SUIT_WHITE, 0},
		card.SuitRankCard{SUIT_WHITE, 0},
		card.SuitRankCard{SUIT_WHITE, 0},
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

func TestStartOfGame(t *testing.T) {
	game := mockGame(t)

	Convey("Given the game has just started", t, func() {
		Convey("It should not be the end of the game", func() {
			So(game.IsFinished(), ShouldBeFalse)
		})
		Convey("The turn phase should be playing or discarding", func() {
			So(game.TurnPhase, ShouldEqual, TURN_PHASE_PLAY_OR_DISCARD)
		})
	})
}

func TestDiscardCard(t *testing.T) {
	game := mockGame(t)

	Convey("Given Steve tries to discard when it's Mick's turn", t, func() {
		game := cloneGame(game)
		_, err := command.CallInCommands("Steve", game, "discard r5",
			game.Commands())
		Convey("It should error", func() {
			So(err, ShouldNotBeNil)
		})
	})
	Convey("Given Mick discards his red 5", t, func() {
		game := cloneGame(game)
		_, err := command.CallInCommands("Mick", game, "discard r5",
			game.Commands())
		Convey("It should not error", func() {
			So(err, ShouldBeNil)
		})
		c := card.SuitRankCard{SUIT_RED, 5}
		Convey("It should take the card from Mick's hand", func() {
			So(len(game.Board.PlayerHands[0]), ShouldEqual, 7)
			So(game.Board.PlayerHands[0].Contains(c), ShouldEqual, 0)
		})
		Convey("It should put the card into the red discard pile", func() {
			So(len(game.Board.DiscardPiles[SUIT_RED]), ShouldEqual, 1)
			So(game.Board.DiscardPiles[SUIT_RED].Contains(c), ShouldEqual, 1)
		})
		Convey("It should change turn phase to draw", func() {
			So(game.TurnPhase, ShouldEqual, TURN_PHASE_DRAW)
		})
	})
}

func TestPlayCard(t *testing.T) {
	game := mockGame(t)

	Convey("Given Mick plays his red 5", t, func() {
		game := cloneGame(game)
		_, err := command.CallInCommands("Mick", game, "play r5",
			game.Commands())
		Convey("It should not error", func() {
			So(err, ShouldBeNil)
		})
		c := card.SuitRankCard{SUIT_RED, 5}
		Convey("It should take the card from Mick's hand", func() {
			So(len(game.Board.PlayerHands[0]), ShouldEqual, 7)
			So(game.Board.PlayerHands[0].Contains(c), ShouldEqual, 0)
		})
		Convey("It should put the card into the Mick's red pile", func() {
			So(len(game.Board.PlayerExpeditions[0][SUIT_RED]), ShouldEqual, 1)
			So(game.Board.PlayerExpeditions[0][SUIT_RED].Contains(c),
				ShouldEqual, 1)
		})
		Convey("It should change turn phase to draw", func() {
			So(game.TurnPhase, ShouldEqual, TURN_PHASE_DRAW)
		})
	})
}

func TestPlaySecondInvestmentCard(t *testing.T) {
	game := mockGame(t)

	// Customise Mick's hand to give extra investment cards
	game.Board.PlayerHands[0] = card.Deck{
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 4},
		card.SuitRankCard{SUIT_RED, 5},
		card.SuitRankCard{SUIT_YELLOW, 0},
		card.SuitRankCard{SUIT_YELLOW, 3},
		card.SuitRankCard{SUIT_WHITE, 10},
	}

	Convey("Given Mick already has a red investment in his red expedition", t,
		func() {
			game.Board.PlayerExpeditions[0][SUIT_RED] = card.Deck{
				card.SuitRankCard{
					Suit: SUIT_RED,
					Rank: 0,
				},
			}
			Convey("Given Mick tries to play a red investment", func() {
				game := cloneGame(game)
				_, err := command.CallInCommands("Mick", game, "play rx",
					game.Commands())
				Convey("It should not error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
}

func TestPlayLowerCard(t *testing.T) {
	game := mockGame(t)

	// Customise Mick's hand to give extra investment cards
	game.Board.PlayerHands[0] = card.Deck{
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 0},
		card.SuitRankCard{SUIT_RED, 4},
		card.SuitRankCard{SUIT_RED, 5},
		card.SuitRankCard{SUIT_YELLOW, 0},
		card.SuitRankCard{SUIT_YELLOW, 3},
		card.SuitRankCard{SUIT_WHITE, 10},
	}

	Convey("Given Mick already has red 3 in his red expedition", t, func() {
		game.Board.PlayerExpeditions[0][SUIT_RED] = card.Deck{
			card.SuitRankCard{
				Suit: SUIT_RED,
				Rank: 3,
			},
		}
		Convey("Given Mick tries to play a red investment", func() {
			game := cloneGame(game)
			_, err := command.CallInCommands("Mick", game, "play rx",
				game.Commands())
			Convey("It should error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestPlayCardNotInHand(t *testing.T) {
	game := mockGame(t)

	Convey("Given Mick does not have red 10", t, func() {
		Convey("Given Mick tries to play red 10", func() {
			game := cloneGame(game)
			expLen := len(game.Board.PlayerExpeditions[0][SUIT_RED])
			_, err := command.CallInCommands("Mick", game, "play r10",
				game.Commands())
			Convey("It should error", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("It should not add the card to the expedition", func() {
				So(len(game.Board.PlayerExpeditions[0][SUIT_RED]), ShouldEqual,
					expLen)
			})
		})
	})
}

func TestDrawCard(t *testing.T) {
	game := mockGame(t)
	if _, err := command.CallInCommands("Mick", game, "discard r5",
		game.Commands()); err != nil {
		t.Fatal(err)
	}

	Convey("Given Mick draws from the draw pile", t, func() {
		game := cloneGame(game)
		topCard, _ := game.Board.DrawPile.Pop()
		initialHandCount := game.Board.PlayerHands[0].Contains(topCard)
		initialDrawCount := game.Board.DrawPile.Contains(topCard)
		_, err := command.CallInCommands("Mick", game, "draw",
			game.Commands())
		Convey("It should not error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Mick should have 8 cards again", func() {
			So(len(game.Board.PlayerHands[0]), ShouldEqual, 8)
		})
		Convey("Mick's hand should have the top card of the deck", func() {
			So(game.Board.PlayerHands[0].Contains(topCard), ShouldEqual,
				initialHandCount+1)
		})
		Convey("The deck should no longer have the top card", func() {
			So(game.Board.DrawPile.Contains(topCard), ShouldEqual,
				initialDrawCount-1)
		})
		Convey("It should be Steve's turn", func() {
			So(game.CurrentlyMoving, ShouldEqual, 1)
		})
		Convey("It should be the play or discard phase", func() {
			So(game.TurnPhase, ShouldEqual, TURN_PHASE_PLAY_OR_DISCARD)
		})
	})
}

func TestTakeCard(t *testing.T) {
	game := mockGame(t)
	if _, err := command.CallInCommands("Mick", game, "discard y3",
		game.Commands()); err != nil {
		t.Fatal(err)
	}
	c := card.SuitRankCard{SUIT_RED, 9}

	Convey("When Mick tries to take from the blue discard pile", t, func() {
		game := cloneGame(game)
		_, err := command.CallInCommands("Mick", game, "take b",
			game.Commands())
		Convey("It should error", func() {
			So(err, ShouldNotBeNil)
		})
	})
	Convey("When there is a red 9 in the discard pile", t, func() {
		game := cloneGame(game)
		game.Board.DiscardPiles[SUIT_RED] = card.Deck{c}
		Convey("When Mick tries to take from the red discard pile", func() {
			game := cloneGame(game)
			_, err := command.CallInCommands("Mick", game, "take r",
				game.Commands())
			Convey("It should not error", func() {
				So(err, ShouldBeNil)
			})
			Convey("The red 9 should be removed from the red discard pile",
				func() {
					So(game.Board.DiscardPiles[SUIT_RED].Contains(c),
						ShouldEqual, 0)
				})
			Convey("The red 9 should be in Mick's hand", func() {
				So(game.Board.PlayerHands[0].Contains(c), ShouldEqual, 1)

			})
			Convey("It should be Steve's turn", func() {
				So(game.CurrentlyMoving, ShouldEqual, 1)
			})
			Convey("It should be the play or discard phase", func() {
				So(game.TurnPhase, ShouldEqual, TURN_PHASE_PLAY_OR_DISCARD)
			})
		})
	})
}

func TestEndOfRound(t *testing.T) {
	game := mockGame(t)
	if _, err := command.CallInCommands("Mick", game, "discard y3",
		game.Commands()); err != nil {
		t.Fatal(err)
	}
	Convey("Given there is one card left in the draw pile", t, func() {
		game := cloneGame(game)
		game.Board.DrawPile, _ = game.Board.DrawPile.PopN(1)
		Convey("Given it is round 0", func() {
			Convey("When Mick draws a card", func() {
				game := cloneGame(game)
				_, err := command.CallInCommands("Mick", game, "draw",
					game.Commands())
				Convey("It should not error", func() {
					So(err, ShouldBeNil)
				})
				Convey("It should be round 1", func() {
					So(game.Round, ShouldEqual, 1)
				})
				Convey("It should not be the end of the game", func() {
					So(game.IsFinished(), ShouldBeFalse)
				})
			})
		})
		Convey("Given it is round 2", func() {
			game := cloneGame(game)
			game.Round = 2
			Convey("When Mick draws a card", func() {
				game := cloneGame(game)
				_, err := command.CallInCommands("Mick", game, "draw",
					game.Commands())
				Convey("It should not error", func() {
					So(err, ShouldBeNil)
				})
				Convey("It should be the end of the game", func() {
					So(game.IsFinished(), ShouldBeTrue)
				})
			})
		})
	})
}

func TestExpeditionScores(t *testing.T) {
	Convey("Given an empty expedition", t, func() {
		expedition := card.Deck{}
		Convey("Score should be 0", func() {
			So(ScoreExpedition(expedition), ShouldEqual, 0)
		})
	})

	Convey("Given an expedition of X", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 0,
			},
		}
		Convey("Score should be -40", func() {
			So(ScoreExpedition(expedition), ShouldEqual, -40)
		})
	})

	Convey("Given an expedition of 3, 4, 5", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 3,
			},
			card.SuitRankCard{
				Rank: 5,
			},
			card.SuitRankCard{
				Rank: 7,
			},
		}
		Convey("Score should be -5", func() {
			So(ScoreExpedition(expedition), ShouldEqual, -5)
		})
	})

	Convey("Given an expedition of X, 5, 7, 10", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 5,
			},
			card.SuitRankCard{
				Rank: 7,
			},
			card.SuitRankCard{
				Rank: 10,
			},
		}
		Convey("Score should be 4", func() {
			So(ScoreExpedition(expedition), ShouldEqual, 4)
		})
	})

	Convey("Given an expedition of X, 2, 3, 4, 5, 6, 7, 10", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 2,
			},
			card.SuitRankCard{
				Rank: 3,
			},
			card.SuitRankCard{
				Rank: 4,
			},
			card.SuitRankCard{
				Rank: 5,
			},
			card.SuitRankCard{
				Rank: 6,
			},
			card.SuitRankCard{
				Rank: 7,
			},
			card.SuitRankCard{
				Rank: 10,
			},
		}
		Convey("Score should be 54", func() {
			So(ScoreExpedition(expedition), ShouldEqual, 54)
		})
	})

	Convey("Given an expedition of X, 2, 3, 4, 5, 6, 7", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 2,
			},
			card.SuitRankCard{
				Rank: 3,
			},
			card.SuitRankCard{
				Rank: 4,
			},
			card.SuitRankCard{
				Rank: 5,
			},
			card.SuitRankCard{
				Rank: 6,
			},
			card.SuitRankCard{
				Rank: 7,
			},
		}
		Convey("Score should be 14", func() {
			So(ScoreExpedition(expedition), ShouldEqual, 14)
		})
	})

	Convey("Given the minimum score expedition of X, X, X", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 0,
			},
		}
		Convey("Score should be -80", func() {
			So(ScoreExpedition(expedition), ShouldEqual, -80)
		})
	})

	Convey("Given the maximum score expedition of X, X, X and 2-10", t, func() {
		expedition := card.Deck{
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 0,
			},
			card.SuitRankCard{
				Rank: 2,
			},
			card.SuitRankCard{
				Rank: 3,
			},
			card.SuitRankCard{
				Rank: 4,
			},
			card.SuitRankCard{
				Rank: 5,
			},
			card.SuitRankCard{
				Rank: 6,
			},
			card.SuitRankCard{
				Rank: 7,
			},
			card.SuitRankCard{
				Rank: 8,
			},
			card.SuitRankCard{
				Rank: 9,
			},
			card.SuitRankCard{
				Rank: 10,
			},
		}
		Convey("Score should be 156", func() {
			So(ScoreExpedition(expedition), ShouldEqual, 156)
		})
	})
}

func TestCannotTakeCardWhenNotTurn(t *testing.T) {
	game := mockGame(t)
	if _, err := command.CallInCommands("Mick", game, "discard y3",
		game.Commands()); err != nil {
		t.Fatal(err)
	}
	if _, err := command.CallInCommands("Steve", game, "take y",
		game.Commands()); err == nil {
		t.Fatal("It allowed Steve to take a turn when it wasn't his turn")
	}
}
