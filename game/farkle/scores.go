package farkle

import (
	"bytes"
	"errors"
	"strconv"
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

func DiceToScoreString(dice []int) (string, error) {
	buf := bytes.NewBufferString("")
	for _, d := range dice {
		if d < 1 || d > 6 {
			return "", errors.New("Can only use numbers 1 to 6 in score string")
		}
		buf.WriteString(strconv.Itoa(d))
	}
	return buf.String(), nil
}

func ScoreStringToDice(scoreString string) ([]int, error) {
	dice := []int{}
	for _, b := range scoreString {
		d := int(b) - int('1') + 1
		if d < 1 || d > 6 {
			return []int{}, errors.New(
				"Can only use numbers 1 to 6 in score string")
		}
		dice = append(dice, d)
	}
	return dice, nil
}

func DiceInDice(search []int, in []int) (isIn bool, remaining []int) {
	searchMap := map[int]int{}
	for _, d := range search {
		searchMap[d]++
	}
	inMap := map[int]int{}
	for _, d := range in {
		inMap[d]++
	}
	for i := 1; i <= 6; i++ {
		if searchMap[i] > inMap[i] {
			return false, in
		}
		for d := 0; d < inMap[i]-searchMap[i]; d++ {
			remaining = append(remaining, i)
		}
	}
	return true, remaining
}

func DiceEquals(a []int, b []int) bool {
	isIn, remaining := DiceInDice(a, b)
	return isIn && len(remaining) == 0
}
