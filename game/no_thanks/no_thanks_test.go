package no_thanks

import (
	"testing"
)

func TestStart(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve"})
	if err != nil {
		t.Error(err)
		return
	}
	if g.Players[0] != "Mick" || g.Players[1] != "Steve" {
		t.Error("Players aren't Mick and Steve, got", g.Players)
		return
	}
	if g.CurrentlyMoving != "Mick" && g.CurrentlyMoving != "Steve" {
		t.Error("Currently moving not set to Mick or Steve, got",
			g.CurrentlyMoving)
		return
	}
}

func TestAllCards(t *testing.T) {
	g := &Game{}
	cards := g.AllCards()
	if len(cards) != 33 {
		t.Error("There weren't 33 cards, got", len(cards))
		return
	}
	if cards[0] != 3 {
		t.Error("Expected the first card to be 3, got", cards[0])
		return
	}
	if cards[32] != 35 {
		t.Error("Expected the thirty third card to be 35, got", cards[32])
		return
	}
}

func TestInitCards(t *testing.T) {
	g := &Game{}
	g.InitCards()
	if len(g.RemainingCards) != 24 {
		t.Error("Expected there to be 24 cards in the stack, got",
			len(g.RemainingCards))
		return
	}
	for _, c := range g.RemainingCards {
		if c < 3 || c > 35 {
			t.Error("Expected cards to be between 3 and 35, got", c)
			return
		}
	}
}

func TestInitPlayerChips(t *testing.T) {
	g := &Game{}
	g.InitPlayerChips()
	for _, p := range g.Players {
		if g.PlayerChips[p] != 11 {
			t.Error("Expected player chips to be 11, got", g.PlayerChips[p])
			return
		}
	}
}
