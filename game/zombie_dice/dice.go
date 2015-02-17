package zombie_dice

import "github.com/Miniand/brdg.me/render"

const (
	Brain = iota
	Shotgun
	Footprints
)

var Colours = []string{
	render.Green,
	render.Yellow,
	render.Red,
}

type Dice struct {
	Colour string
	Faces  []int
}

func (d Dice) Roll() int {
	return d.Faces[rnd.Int()%len(d.Faces)]
}

func (d Dice) RollResult() DiceResult {
	return DiceResult{d, d.Roll()}
}

func RollDice(dice []Dice) DiceResultList {
	drl := make(DiceResultList, len(dice))
	for i, die := range dice {
		drl[i] = die.RollResult()
	}
	return drl
}

var GreenDice = Dice{
	render.Green,
	[]int{
		Brain,
		Brain,
		Brain,
		Footprints,
		Footprints,
		Shotgun,
	},
}

var YellowDice = Dice{
	render.Yellow,
	[]int{
		Brain,
		Brain,
		Footprints,
		Footprints,
		Shotgun,
		Shotgun,
	},
}

var RedDice = Dice{
	render.Red,
	[]int{
		Brain,
		Footprints,
		Footprints,
		Shotgun,
		Shotgun,
		Shotgun,
	},
}

type DiceResult struct {
	Dice
	Face int
}

type DiceResultList []DiceResult

var ColourOrder = map[string]int{
	render.Green:  0,
	render.Yellow: 1,
	render.Red:    2,
}

func (drl DiceResultList) Len() int      { return len(drl) }
func (drl DiceResultList) Swap(i, j int) { drl[i], drl[j] = drl[j], drl[i] }
func (drl DiceResultList) Less(i, j int) bool {
	if drl[i].Face != drl[j].Face {
		return drl[i].Face < drl[j].Face
	}
	return ColourOrder[drl[i].Colour] < ColourOrder[drl[j].Colour]
}

func (drl DiceResultList) Dice() []Dice {
	dice := make([]Dice, len(drl))
	for i, dr := range drl {
		dice[i] = dr.Dice
	}
	return dice
}

func AllDice() []Dice {
	return []Dice{
		GreenDice,
		GreenDice,
		GreenDice,
		GreenDice,
		GreenDice,
		GreenDice,
		YellowDice,
		YellowDice,
		YellowDice,
		YellowDice,
		RedDice,
		RedDice,
		RedDice,
	}
}
