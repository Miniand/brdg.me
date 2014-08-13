package sushizock

import (
	"math/rand"
	"time"
)

const (
	DiceSushi = iota
	DiceBlueChopsticks
	DiceBones
	DiceRedChopsticks
)

var DiceFaces = []int{
	DiceSushi,
	DiceSushi,
	DiceBones,
	DiceBones,
	DiceBlueChopsticks,
	DiceRedChopsticks,
}

var DiceText = map[int]string{
	DiceSushi:          `{{c "blue"}}Θ{{_c}}`,
	DiceBlueChopsticks: `{{c "blue"}}X{{_c}}`,
	DiceBones:          `{{c "red"}}¥{{_c}}`,
	DiceRedChopsticks:  `{{c "red"}}X{{_c}}`,
}

func RollDie() int {
	return DiceFaces[rand.New(rand.NewSource(time.Now().UnixNano())).Int()%
		len(DiceFaces)]
}

func RollDice(n int) []int {
	dice := make([]int, n)
	for i := 0; i < n; i++ {
		dice[i] = RollDie()
	}
	return dice
}

func DiceCounts(dice []int) map[int]int {
	counts := map[int]int{}
	for _, d := range dice {
		counts[d] += 1
	}
	return counts
}

func DiceStrings(dice []int) []string {
	strs := make([]string, len(dice))
	for i, d := range dice {
		strs[i] = DiceText[d]
	}
	return strs
}
