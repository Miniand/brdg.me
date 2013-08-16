package die

import (
	"bytes"
	"errors"
	"strconv"
)

func Render(value int) (output string) {
	switch value {
	case 1:
		output = "⚀"
	case 2:
		output = "⚁"
	case 3:
		output = "⚂"
	case 4:
		output = "⚃"
	case 5:
		output = "⚄"
	case 6:
		output = "⚅"
	default:
		panic("A die value is between 1 and 6")
	}
	return
}

func DiceToValueString(dice []int) (string, error) {
	buf := bytes.NewBufferString("")
	for _, d := range dice {
		if d < 1 || d > 6 {
			return "", errors.New("Can only use numbers 1 to 6 in score string")
		}
		buf.WriteString(strconv.Itoa(d))
	}
	return buf.String(), nil
}

func ValueStringToDice(scoreString string) ([]int, error) {
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
