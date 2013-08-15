package farkle

import (
	"testing"
)

func TestDiceInDice(t *testing.T) {
	isIn, remaining := DiceInDice([]int{1, 1, 1}, []int{2, 4, 1, 3, 1, 1})
	if !isIn {
		t.Fatal("Expected to be in")
	}
	if !DiceEquals(remaining, []int{2, 4, 3}) {
		t.Fatal(remaining)
	}
	isIn, remaining = DiceInDice([]int{5}, []int{1, 4, 5, 5, 5, 3})
	if !isIn {
		t.Fatal("Expected to be in")
	}
	if !DiceEquals(remaining, []int{1, 4, 5, 5, 3}) {
		t.Fatal(remaining)
	}
	isIn, remaining = DiceInDice([]int{6}, []int{1, 4, 5, 5, 5, 3})
	if isIn {
		t.Fatal("Expected to not be in")
	}
}
