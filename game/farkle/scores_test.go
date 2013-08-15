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

func TestDiceToScoreString(t *testing.T) {
	str, err := DiceToScoreString([]int{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	if str != "1234" {
		t.Fatal("Expected string to be 1234, got:", str)
	}
	str, err = DiceToScoreString([]int{0})
	if err == nil {
		t.Fatal("Expected an error")
	}
	str, err = DiceToScoreString([]int{7})
	if err == nil {
		t.Fatal("Expected an error")
	}
}

func TestScoreStringToDice(t *testing.T) {
	dice, err := ScoreStringToDice("11153")
	if err != nil {
		t.Fatal(err)
	}
	expect := []int{1, 1, 1, 5, 3}
	if !DiceEquals(dice, expect) {
		t.Fatal("Expected dice to be", expect, "but got", dice)
	}
}
