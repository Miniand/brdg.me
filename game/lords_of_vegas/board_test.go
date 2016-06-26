package lords_of_vegas

import (
	"testing"

	"github.com/bmizerany/assert"
)

var expectedBuildPriceForDice = map[int]int{
	1: 8,
	2: 6,
	3: 9,
	4: 12,
	5: 15,
	6: 20,
}

func TestBoardSpaces(t *testing.T) {
	for _, bs := range BoardSpaces {
		// StartingMoney should be 10-Dice
		assert.Equal(t, 10-bs.Dice, bs.StartingMoney)
		// Check build prices based on dice
		assert.Equal(t, expectedBuildPriceForDice[bs.Dice], bs.BuildPrice)
	}
}
