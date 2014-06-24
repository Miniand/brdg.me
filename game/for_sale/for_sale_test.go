package for_sale

import (
	"reflect"
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
)

const (
	Mick  = "Mick"
	Steve = "Steve"
	BJ    = "BJ"
)

func TestFullGame(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{Mick, Steve, BJ}); err != nil {
		t.Fatal(err)
	}
	// Set the state of the game to sorted decks
	_, g.BuildingDeck = BuildingDeck().PopN(2)
	_, g.ChequeDeck = ChequeDeck().PopN(2)
	g.OpenCards, g.BuildingDeck = g.BuildingDeck.PopN(3)
	if !reflect.DeepEqual([]string{Mick}, g.WhoseTurn()) {
		t.Fatal("It's not Mick's turn.")
	}
	if _, err := command.CallInCommands(
		Mick, g, "bid 3", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{Steve}, g.WhoseTurn()) {
		t.Fatal("It's not Steve's turn.")
	}
	if _, err := command.CallInCommands(
		Steve, g, "bid 3", g.Commands()); err == nil {
		t.Fatal("It let Steve make the same bid.")
	}
	if _, err := command.CallInCommands(
		Steve, g, "bid 4", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{BJ}, g.WhoseTurn()) {
		t.Fatal("It's not BJ's turn.")
	}
	if _, err := command.CallInCommands(
		BJ, g, "pass", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 17},
		card.SuitRankCard{Rank: 18},
	}, g.OpenCards) {
		t.Fatalf("Open cards remaining aren't a 17 and 18, got:\n\n%#v",
			g.OpenCards)
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 16},
	}, g.Hands[2]) {
		t.Fatalf("BJ's hand is not a 16, got:\n\n%#v", g.Hands[2])
	}
	if g.Chips[2] != 15 {
		t.Fatal("BJ doesn't have 15 chips remaining.")
	}
	if !reflect.DeepEqual([]string{Mick}, g.WhoseTurn()) {
		t.Fatal("It's not Mick's turn.")
	}
	if _, err := command.CallInCommands(
		Mick, g, "pass", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if g.Chips[0] != 14 {
		t.Fatal("Mick doesn't have 14 chips remaining.")
	}
	if g.Chips[1] != 11 {
		t.Fatal("Steve doesn't have 11 chips remaining.")
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 17},
	}, g.Hands[0]) {
		t.Fatalf("BJ's hand is not a 17, got:\n\n%#v", g.Hands[0])
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 18},
	}, g.Hands[1]) {
		t.Fatalf("BJ's hand is not a 18, got:\n\n%#v", g.Hands[1])
	}
	if !reflect.DeepEqual([]string{Steve}, g.WhoseTurn()) {
		t.Fatal("It's not Steve's turn.")
	}
	// End the buying phase early and shorten the selling phase.
	g.BuildingDeck = card.Deck{}
	_, g.ChequeDeck = g.ChequeDeck.PopN(15)
	g.OpenCards = card.Deck{}
	g.StartRound()
	if !reflect.DeepEqual([]string{Mick, Steve, BJ}, g.WhoseTurn()) {
		t.Fatal("It's not everyone's turn.")
	}
	if _, err := command.CallInCommands(
		BJ, g, "play 18", g.Commands()); err == nil {
		t.Fatal("It let BJ play an 18 even though it's not in his hand")
	}
	if _, err := command.CallInCommands(
		BJ, g, "play 16", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{Mick, Steve}, g.WhoseTurn()) {
		t.Fatal("It's not Mick and Steve's turn.")
	}
	if _, err := command.CallInCommands(
		Steve, g, "play 18", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{Mick}, g.WhoseTurn()) {
		t.Fatal("It's not Mick's turn.")
	}
	if _, err := command.CallInCommands(
		Mick, g, "play 17", g.Commands()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 0},
	}, g.Cheques[0]) {
		t.Fatal("Mick didn't get the 0 cheque")
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 3},
	}, g.Cheques[1]) {
		t.Fatal("Steve didn't get the 3 cheque")
	}
	if !reflect.DeepEqual(card.Deck{
		card.SuitRankCard{Rank: 0},
	}, g.Cheques[2]) {
		t.Fatal("BJ didn't get the 0 cheque")
	}
	if !g.IsFinished() {
		t.Fatal("The game isn't finished")
	}
	if !reflect.DeepEqual([]string{Steve}, g.Winners()) {
		t.Fatal("Steve was not the winner.")
	}
}
