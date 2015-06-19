package farkle

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/die"
)

type Score struct {
	Dice  []int
	Value int
}

func Scores() []Score {
	return []Score{
		Score{[]int{5}, 50},
		Score{[]int{1}, 100},
		Score{[]int{2, 2, 2}, 200},
		Score{[]int{3, 3, 3}, 300},
		Score{[]int{4, 4, 4}, 400},
		Score{[]int{5, 5, 5}, 500},
		Score{[]int{6, 6, 6}, 600},
		Score{[]int{1, 1, 1}, 1000},
	}
}

func (s Score) ValueString() string {
	valueString, err := die.DiceToValueString(s.Dice)
	if err != nil {
		panic(err.Error())
	}
	return valueString
}

func (s Score) Description() string {
	return fmt.Sprintf("%s (%d points)", s.ValueString(), s.Value)
}

func ScoreStrings() (scoreStrings []string) {
	for _, s := range Scores() {
		valueString, err := die.DiceToValueString(s.Dice)
		if err != nil {
			panic(err.Error())
		}
		scoreStrings = append(scoreStrings, fmt.Sprintf("%s (%d points)",
			valueString, s.Value))
	}
	return
}

func AvailableScores(dice []int) (available map[string]Score) {
	available = map[string]Score{}
	for _, s := range Scores() {
		isIn, _ := die.DiceInDice(s.Dice, dice)
		if isIn {
			available[s.ValueString()] = s
		}
	}
	return
}
